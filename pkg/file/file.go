package file

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
)

//对文件处理措施

func GetSize(f multipart.File) (int, error) {
	content, err := io.ReadAll(f)
	return len(content), err
}

func GetExt(fileName string) string {
	return path.Ext(fileName)
}

func CheckNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist {
		if err := MkDir(src); err != nil {
			return err
		}
	}
	return nil
}
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return f, err
}

func MustOpen(filename, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	src := dir + "/" + filePath
	perm := CheckPermission(src)
	if perm {
		return nil, fmt.Errorf("file.CheckPermission err :%v", err)
	}
	err = IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("IsNotExistMkDir err :%v", err)
	}
	f, err := Open(src+filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("Open err :%v", err)
	}
	return f, nil
}
