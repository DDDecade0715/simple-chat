package helper

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"os"
	"strconv"
	"unicode/utf8"
)

func InArray(needle interface{}, hystack interface{}) bool {
	switch key := needle.(type) {
	case string:
		for _, item := range hystack.([]string) {
			if key == item {
				return true
			}
		}
	case int:
		for _, item := range hystack.([]int) {
			if key == item {
				return true
			}
		}
	case int64:
		for _, item := range hystack.([]int64) {
			if key == item {
				return true
			}
		}
	default:
		return false
	}
	return false
}

func Md5Encrypt(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	return hex.EncodeToString(m.Sum(nil))
}

func MbStrLen(str string) int {
	return utf8.RuneCountInString(str)
}

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
		fi, err := os.Stat(tmpFilePath + fileName + "_" + iStr)
		if err == nil {
			chunkSize += fi.Size()
		}
	}
	if chunkSize == int64(fileSize) {
		return true
	}
	return false
}
