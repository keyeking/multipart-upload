package core

import (
	"multipart-upload/global"
	"multipart-upload/server"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	global.Router = gin.Default()
	global.Router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowWebSockets:  true,
		MaxAge:           12 * time.Hour,
	}))
	global.Router.LoadHTMLGlob("template/*")
	{
		global.Router.GET("/", server.Index)
		global.Router.GET("/list", server.List)
		global.Router.POST("/startUpload", server.StartUpload)
		global.Router.POST("/upload", server.Upload)
		global.Router.POST("/endtUpload", server.EndUpload)
	}
}
