package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//获取单个文件的大小
func GetFileSize(path string) (int64, bool) {
	fileInfo, err := os.Stat(path)
	if err == nil {
		fileSize := fileInfo.Size() //获取size
		return fileSize, true
	}
	return 0, false
}

func GetFileList(path string) []string {
	result := []string{}
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, "xlsm") || strings.HasSuffix(path, "xlsx") {

			result = append(result, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
		return result
	}
	return result
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

// 递归删除目录
func DelDir(path string) bool {
	if err := os.RemoveAll(path); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func MkDir(path string) bool {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
