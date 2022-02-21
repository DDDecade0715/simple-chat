package file

import (
	"errors"
	"os"
	"strconv"
)

//CreateFile 创建文件
func CreateFile(filePath string) (bool, error) {
	fileBool, err := FileExists(filePath)
	if fileBool && err == nil {
		return true, errors.New("文件以存在")
	} else {
		newFile, err := os.Create(filePath)
		defer newFile.Close()
		if err != nil {
			return false, errors.New("创建文件失败")
		}
	}
	return true, nil
}

//FileExists 判断文件或文件夹是否存在
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return false, err
}

//IsFnish 判断是否完成  根据现有文件的大小 与 上传文件大小进行匹配
func IsFnish(fileName string, chunkTotal, fileSize int, tmpFilePath string) bool {
	var chunkSize int64
	for i := 0; i < chunkTotal; i++ {
		iStr := strconv.Itoa(i)
		// 分片大小获取
		fi, err := os.Stat(tmpFilePath + "/" + fileName + "_" + iStr)
		if err == nil {
			chunkSize += fi.Size()
		}
	}
	if chunkSize == int64(fileSize) {
		return true
	}
	return false
}
