package export

import "gogin/pkg/setting"

//关于excel的各种路径和url

func GetExcelFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetExcelPath() + name
}

func GetExcelFullSavePath() string {
	return setting.AppSetting.RuntimeRootPath + GetExcelPath()
}

func GetExcelPath() string {
	return setting.AppSetting.ExportExcelSavePath
}
