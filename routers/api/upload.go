package api

import (
	"github.com/gin-gonic/gin"
	"gogin/pkg/app"
	"gogin/pkg/logging"
	"gogin/pkg/mistakeMsg"
	"gogin/pkg/upload"
	"net/http"
)

func UploadImage(c *gin.Context) {
	appG := app.Response{c}
	code := mistakeMsg.SUCCESS
	data := make(map[string]string)

	file, image, err := c.Request.FormFile("image")

	if err != nil {
		logging.Warn(err)
		code = mistakeMsg.ERROR
		appG.ResponseJson(http.StatusOK, code, data)
		return
	}
	if image == nil {
		code = mistakeMsg.INVALID_PARAMS
		appG.ResponseJson(http.StatusOK, code, data)
	}

	imageName := upload.GetImageName(image.Filename) //md5编码
	fullPath := upload.GetImageFullPath()
	savePath := upload.GetImagePath()

	src := fullPath + imageName
	if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
		code = mistakeMsg.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
		appG.ResponseJson(http.StatusOK, code, data)
		return
	}
	err = upload.CheckImage(fullPath)
	if err != nil {
		logging.Warn(err)
		code = mistakeMsg.ERROR_UPLOAD_CHECK_IMAGE_FAIL
		appG.ResponseJson(http.StatusOK, code, data)
		return
	}
	if err := c.SaveUploadedFile(image, src); err != nil {
		logging.Warn(err)
		code = mistakeMsg.ERROR_UPLOAD_SAVE_IMAGE_FAIL
		appG.ResponseJson(http.StatusOK, code, data)
		return
	}

	data["image_url"] = upload.GetImageFullUrl(imageName)
	data["image_save_url"] = savePath + imageName
	appG.ResponseJson(http.StatusOK, code, data)
}
