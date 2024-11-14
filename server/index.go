package server

import (
	"fmt"
	"multipart-upload/global"
	"multipart-upload/models"
	"multipart-upload/utils"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "分片上传",
	})
}
func List(c *gin.Context) {
	id := c.Query("userId")
	var fileList []models.FileModel
	if err := global.Db.Model(&models.FileModel{}).Preload("FileChunkModel").Where("user_id=?", id).Find(&fileList).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": fileList,
	})
}

// 开始上传
func StartUpload(c *gin.Context) {
	//获取文件信息
	var f models.FileModel
	if err := c.ShouldBindBodyWithJSON(&f); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": err,
		})
		return
	}

	// 判断数据库中是否存在，断点续传情况
	//主要是判断文件名、用户id、切片数量、大小、类型是否一致
	var f2 models.FileModel
	if err := global.Db.Model(&models.FileModel{}).Preload("FileChunkModel", func(db *gorm.DB) *gorm.DB {
		return db.Order("file_chunk_model.index ASC")
	}).Where("user_id=? and source_file_name=? and type = ? and size = ? and slice_count = ?", f.UserId, f.FileName, f.Type, f.Size, f.SliceCount).Find(&f2).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": err,
		})
		return
	}
	// 文件存在，判断切片缺少情况
	if len(f2.FileChunkModel) != 0 {
		var missingChunks = make([]int, 0)
		var Chunks = make([]int, 0)
		for _, v := range f2.FileChunkModel {
			Chunks = append(Chunks, v.Index)
		}
		for i := 1; i <= f2.SliceCount; i++ {
			if b := utils.ContainsInt(Chunks, i); !b {
				missingChunks = append(missingChunks, i)
			}
		}

		process := fmt.Sprintf("保留两位小数后的结果: %.2f%%", float64(len(f2.FileChunkModel))/float64(f2.SliceCount)*100)
		c.JSON(http.StatusOK, gin.H{
			"code": 2000, //标志断点续传
			"data": f2,
			"list": missingChunks, //缺少的切片
			"msg":  fmt.Sprintf("文件已存在，已上传%v", process),
		})
	} else {
		f.SourceFileName = f.FileName
		// 添加时间戳，防止重名文件
		fileName := fmt.Sprintf("%s_%v%s", strings.TrimSuffix(f.FileName, filepath.Ext(f.FileName)), time.Now().UnixNano(), filepath.Ext(f.FileName))
		f.FileName = fileName
		//寻找是否存在文件目录
		var sourcePath string
		var b bool
		if b, sourcePath = utils.IsExist(f.FileName); !b {
			// 创建目录
			var err error
			if sourcePath, err = utils.MkMultiDir(f.FileName); err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": 201,
					"msg":  err,
				})
				return
			}
		}
		f.FilePath = sourcePath
		if err := global.Db.Model(&models.FileModel{}).Create(&f).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 201,
				"msg":  err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": f,
			"msg":  fmt.Sprintf("文件将保存在%s", sourcePath),
		})
	}
}

// 正在上传
func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	sourcePath := c.PostForm("sourcePath")
	fileName := c.PostForm("fileName")
	index := c.PostForm("index")
	id := c.PostForm("id")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 201,
			"msg":  err,
		})
		return
	}
	pathFile := filepath.Join(sourcePath, fmt.Sprintf("%v-%s.temp", index, fileName))
	// 上传文件到指定的目录
	if err := c.SaveUploadedFile(file, pathFile); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 201,
			"msg":  err,
		})
		return
	}
	i, _ := strconv.Atoi(index)
	var f = models.FileChunkModel{
		FileName:    fmt.Sprintf("%v-%s.temp", i, fileName),
		FilePath:    pathFile,
		Index:       i,
		FileModelId: id,
	}
	if err := global.Db.Model(&models.FileChunkModel{}).Create(&f).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 201,
			"msg":  err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": "upload",
	})
}

// 上传结束
func EndUpload(c *gin.Context) {
	//获取文件信息
	var f models.FileModel
	if err := c.ShouldBindBodyWithJSON(&f); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": err,
		})
		return
	}
	var file models.FileModel
	// 判断数据库中是否存在该文件
	result := global.Db.Model(&models.FileModel{}).Preload("FileChunkModel", func(db *gorm.DB) *gorm.DB {
		return db.Order("file_chunk_model.index ASC")
	}).Where(f).Find(&file)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 201,
			"msg":  result.Error,
		})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 201,
			"msg":  "文件不存在",
		})
		return
	}
	// 判断文件切片列表是否保存完整
	if len(file.FileChunkModel) != f.SliceCount {
		c.JSON(http.StatusOK, gin.H{
			"code": 201,
			"msg":  "文件切片列表不完整",
		})
		return
	}
	// 合并文件
	if err := utils.MergeFile(file); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 201,
			"msg":  "文件合并失败",
		})
		return
	}
	// 删除数据库中的切片信息
	if err := global.Db.Where("file_model_id = ?", file.Id).Unscoped().Delete(&models.FileChunkModel{}).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 201,
			"msg":  err,
		})
		return
	}
	if err := utils.CleanupTempFiles(file.FilePath); err != nil {
		fmt.Println("清理失败")
	}
	c.JSON(http.StatusOK, gin.H{
		"code": file,
	})
}
