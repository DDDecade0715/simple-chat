package messageService

import (
	"errors"
	"gin-derived/api/models"
	"gin-derived/global"
	"gin-derived/pkg/app/response"
	"gin-derived/pkg/file"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"sync"
	"time"
)

var lock sync.WaitGroup

//MergeChunk 合并文件
func MergeChunk(c *gin.Context) (int, error) {
	// 上传文件
	_, err := ChunkUpload(c)
	if err != nil {
		return 0, errors.New("上传失败")
	}
	// 分片总数
	chunkTotal := c.Request.FormValue("chunktotal")
	// 文件总大小
	fileSize := c.Request.FormValue("filesize")
	// 获取上传文件
	_, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		return 0, errors.New("上传文件错误")
	}
	// 分片序号与分片总数相等 则合并文件
	total, _ := strconv.Atoi(chunkTotal)
	size, _ := strconv.Atoi(fileSize)
	// 上传总数
	totalLen := 0
	//拼接文件夹
	savePath := global.GCONFIG.App.UploadPath
	realPath := "/" + time.Now().Format("20060102")
	//按日期分组
	savePath += realPath
	// 最后一个分片上传时,进行分片合并成一个文件
	if file.IsFnish(fileHeader.Filename, total, size, savePath) {
		// 新文件创建
		filePath := savePath + "/" + fileHeader.Filename
		fileBool, err := file.CreateFile(filePath)
		if !fileBool {
			return 0, err
		}
		// 读取文件片段 进行合并
		for i := 0; i < total; i++ {
			lock.Add(1)
			go MergeFile(i, fileHeader.Filename, filePath)
		}
		lock.Wait()

		chatId, ok := c.GetPostForm("chat_id")
		if !ok {
			response.FailWithMessage("fail", c)
			return 0, errors.New("没上传chat_id")
		}

		value := &models.ChatImage{
			MessageId: chatId,
			Url:       global.GCONFIG.App.UploadUrl + realPath + "/" + fileHeader.Filename,
			FileSize:  int64(size),
			FileName:  fileHeader.Filename,
		}
		err = models.CreateChatImage(value)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return 0, errors.New("保存错误")
		}
	}
	return totalLen, nil
}

//MergeFile 合并切片文件
func MergeFile(i int, fileName, filePath string) {
	// 打开之前上传文件
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer file.Close()
	if err != nil {
		global.GLOG.Error("打开之前上传文件不存在")
	}
	//拼接文件夹
	savePath := global.GCONFIG.App.UploadPath
	realPath := "/" + time.Now().Format("20060102")
	//按日期分组
	savePath += realPath
	// 分片大小获取
	fi, _ := os.Stat(savePath + "/" + fileName + "_0")
	chunkSize := fi.Size()
	// 设置文件写入偏移量
	file.Seek(chunkSize*int64(i), 0)
	iSize := strconv.Itoa(i)
	chunkFilePath := savePath + "/" + fileName + "_" + iSize
	chunkFileObj, err := os.Open(chunkFilePath)
	defer chunkFileObj.Close()
	if err != nil {
		global.GLOG.Error("打开分片文件失败")
	}

	// 上传总数
	totalLen := 0
	// 写入数据
	data := make([]byte, 1024, 1024)
	for {
		tal, err := chunkFileObj.Read(data)
		if err == io.EOF {
			// 删除文件 需要先关闭改文件
			chunkFileObj.Close()
			err := os.Remove(chunkFilePath)
			if err != nil {
				global.GLOG.Error("临时记录文件删除失败", err)
			}
			global.GLOG.Info("文件复制完毕")
			break
		}
		len, err := file.Write(data[:tal])
		if err != nil {
			global.GLOG.Error("文件上传失败")
		}
		totalLen += len
	}
	lock.Done()
	//return totalLen,nil
}

//ChunkUpload 分片上传
func ChunkUpload(c *gin.Context) (int, error) {
	// 分片序号
	chunkIndex := c.Request.FormValue("chunkindex")
	// 获取上传文件
	upFile, fileHeader, err := c.Request.FormFile("file")

	if err != nil {
		return 0, errors.New("上传文件错误")
	}

	// 新文件创建
	savePath := global.GCONFIG.App.UploadPath
	realPath := "/" + time.Now().Format("20060102")
	//按日期分组
	savePath += realPath
	//创建文件夹
	_, osErr := os.Stat(savePath)
	if os.IsNotExist(osErr) {
		err = os.MkdirAll(savePath, os.ModePerm)
		if err != nil {
			return 0, errors.New("failed to create save directory.")
		}
	}
	//判断权限
	if os.IsPermission(osErr) {
		return 0, errors.New("insufficient file permissions.")
	}

	filePath := savePath + "/" + fileHeader.Filename + "_" + chunkIndex
	fileBool, err := file.CreateFile(filePath)
	if !fileBool {
		return 0, err
	}
	// 获取现在文件大小
	fi, _ := os.Stat(filePath)
	// 判断文件是否传输完成
	if fi.Size() == fileHeader.Size {
		return 0, errors.New("文件已存在, 不继续上传")
	}
	start := strconv.Itoa(int(fi.Size()))

	// 进行断点上传
	// 打开之前上传文件
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer file.Close()
	if err != nil {
		return 0, errors.New("打开之前上传文件不存在")
	}

	// 将数据写入文件
	count, _ := strconv.ParseInt(start, 10, 64)
	total, err := UploadFile(upFile, count, file, count)
	return total, err
}

// 上传文件
func UploadFile(upfile multipart.File, upSeek int64, file *os.File, fSeek int64) (int, error) {
	// 上传文件大小记录
	fileSzie := 0
	// 设置上传偏移量
	upfile.Seek(upSeek, 0)
	// 设置文件偏移量
	file.Seek(fSeek, 0)
	data := make([]byte, 1024, 1024)
	for {
		total, err := upfile.Read(data)
		if err == io.EOF {
			//fmt.Println("文件复制完毕")
			break
		}
		len, err := file.Write(data[:total])
		if err != nil {
			return 0, errors.New("文件上传失败")
		}
		// 记录上传长度
		fileSzie += len
	}
	return fileSzie, nil
}
