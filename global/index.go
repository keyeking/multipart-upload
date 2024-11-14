package global

import (
	"multipart-upload/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	Router *gin.Engine
	Config *models.Config
	Db     *gorm.DB
)
