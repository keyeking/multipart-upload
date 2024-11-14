package core

import (
	"fmt"
	"multipart-upload/global"
	"multipart-upload/models"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitGorm /*初始化数据库*/
func InitGorm() {
	dsn := global.Config.Gorm.Dsn()
	var mysqlLogger logger.Interface
	// 开发环境显示的sql
	mysqlLogger = logger.Default.LogMode(logger.Info)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:          mysqlLogger,
		CreateBatchSize: 100,
		PrepareStmt:     true,
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		DisableForeignKeyConstraintWhenMigrating: false,
	})
	if err != nil {
		panic(fmt.Sprintf("数据库连接失败: %s", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)               // 设置连接池中的最大闲置连接数
	sqlDB.SetMaxOpenConns(100)              // 设置连接池最大连接数
	sqlDB.SetConnMaxLifetime(time.Hour * 4) // 设置连接池最大生存时间，不能超过mysql的wait_timeout
	global.Db = db
	if global.Config.Gorm.IsMigrate {
		Migrate()
	}
}

func Migrate() {
	fmt.Println("12345")
	err := global.Db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&models.FileModel{},
		&models.FileChunkModel{},
	)
	if err != nil {
		panic(fmt.Sprintf("迁移数据库失败：%s", err))
	}
	fmt.Println("迁移数据库完成")
}
