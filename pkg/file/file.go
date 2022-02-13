package file

import (
	"errors"
	"fmt"
	"gin-derived/global"
	"gin-derived/pkg/helper"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"sync"
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
	// 最后一个分片上传时,进行分片合并成一个文件
	if helper.IsFnish(fileHeader.Filename, total, size, global.GCONFIG.App.UploadPath) {
		// 新文件创建
		filePath := "./" + global.GCONFIG.App.UploadPath + "/" + fileHeader.Filename
		fileBool, err := helper.CreateFile(filePath)
		if !fileBool {
			return 0, err
		}
		// 读取文件片段 进行合并
		for i := 0; i < total; i++ {
			lock.Add(1)
			go MergeFile(i, fileHeader.Filename, filePath)
		}
		lock.Wait()
	}
	return totalLen, nil
}

//MergeFile 合并切片文件
func MergeFile(i int, fileName, filePath string) {
	// 打开之前上传文件
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer file.Close()
	if err != nil {
		log.Fatal("打开之前上传文件不存在")
		//return 0, errors.New("打开之前上传文件不存在")
	}
	// 分片大小获取
	fi, _ := os.Stat(global.GCONFIG.App.UploadPath + "/" + fileName + "_0")
	chunkSize := fi.Size()
	// 设置文件写入偏移量
	file.Seek(chunkSize*int64(i), 0)
	iSize := strconv.Itoa(i)
	chunkFilePath := global.GCONFIG.App.UploadPath + "/" + fileName + "_" + iSize
	fmt.Printf("分片路径:", chunkFilePath)
	chunkFileObj, err := os.Open(chunkFilePath)
	defer chunkFileObj.Close()
	if err != nil {
		log.Fatal("打开分片文件失败")
		//return 0, errors.New("打开分片文件失败")
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
				fmt.Println("临时记录文件删除失败", err)
			}
			fmt.Println("文件复制完毕")
			break
		}
		len, err := file.Write(data[:tal])
		if err != nil {
			log.Fatal("文件上传失败")
			//return 0, errors.New("文件上传失败")
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
	filePath := global.GCONFIG.App.UploadPath + "/" + fileHeader.Filename + "_" + chunkIndex
	fileBool, err := helper.CreateFile(filePath)
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
