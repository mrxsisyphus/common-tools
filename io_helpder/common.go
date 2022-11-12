package io_helpder

import "os"

// PathExists 判断给定路径(文件/文件夹)是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		//如果返回的错误为nil,说明文件或文件夹存在
		return true, nil
	}
	if os.IsNotExist(err) {
		//如果返回的错误类型使用os.IsNotExist()判断为true,说明文件或文件夹不存在
		return false, nil
	}
	//如果返回的错误为其它类型,则不确定是否在存在
	return false, err
}

// IsDir 判断给定路径是否是文件夹
func IsDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false //文件不存在
	}
	return stat.IsDir()
}

// IsFile 判断给定路径是否是文件
func IsFile(path string) bool {
	return !IsDir(path)
}
