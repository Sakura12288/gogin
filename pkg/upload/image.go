package upload

import (
	"fmt"
	"gogin/pkg/file"
	"gogin/pkg/logging"
	"gogin/pkg/setting"
	"gogin/pkg/util"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

//上传图片

func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}

func GetImageFullUrl(name string) string {
	return setting.AppSetting.ImagePrefixUrl + "/" + GetImagePath() + name
}

func GetImageName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)
	return fileName + ext
}

func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

func CheckImageExt(fileName string) bool {
	ext := file.GetExt(fileName)
	for _, v := range setting.AppSetting.ImageAllowExts {
		if strings.ToUpper(v) == strings.ToUpper(ext) {
			return true
		}
	}
	return false
}

func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		log.Println(err)
		logging.Warn(err)
		return false
	}
	return size <= setting.AppSetting.ImageMaxsize
}

func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("CheckImage os.Getwd err : %v", err)
	}
	err = file.IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("CheckImage CheckImage err : %v", err)
	}
	perm := file.CheckPermission(src)
	if perm {
		return fmt.Errorf("CheckImage file.CheckPermission denied src: %s", src)
	}
	return nil
}
