package util

import (
	"os"
	"strconv"
)

// ByteToMb byte转KB
func ByteToMb(b int64, prec int) (KB string) {
	factor := float64(1024)
	fb := float64(b)
	KB = strconv.FormatFloat(fb/factor, 'f', prec, 64)
	return
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

// 载入配置文件
func ConfigLoader(path string) {

}
