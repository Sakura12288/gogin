package routers

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "gogin/docs"
	"gogin/middlewares/jwt"
	"gogin/pkg/export"
	"gogin/pkg/qrcode"
	"gogin/pkg/setting"
	"gogin/pkg/upload"
	"gogin/routers/api"
	v1 "gogin/routers/api/v1"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	gin.SetMode(setting.ServerSetting.RunMode)

	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	r.StaticFS("/export", http.Dir(export.GetExcelFullSavePath()))
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullSavePath()))

	r.POST("/authHS256", api.GetAuthUsingHS256) //获取token应该设置在全局
	r.POST("/authRS256", api.GetAuthUsingRS256)
	r.POST("/upload", api.UploadImage)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JwtRS256()) //获取文件时进行验证

	{
		apiv1.GET("/tags", v1.GetTags)
		apiv1.POST("/tags", v1.AddTag)
		apiv1.PUT("/tags/:id", v1.EditTag)
		apiv1.DELETE("/tags/:id", v1.DeleteTag)
		r.POST("/tags/export", v1.ExportTags)
		r.POST("tags/import", v1.ImportTags)

		apiv1.GET("/articles", v1.GetArticles)
		apiv1.GET("/articles/:id", v1.GetArticle)
		apiv1.POST("/articles", v1.AddArticles)
		apiv1.PUT("/articles/:id", v1.EditArticles)
		apiv1.DELETE("/articles/:id", v1.DeleteArticles)
		apiv1.POST("/articles/poster/generate", v1.GenerateArticlePoster)

	}
	return r
}
