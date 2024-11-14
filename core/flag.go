package core

import (
	"flag"
	"multipart-upload/global"
)

func InitFlag() {
	flag.BoolVar(&global.Config.Gorm.IsMigrate, "migrate", false, "是否进行数据库迁移")
	flag.Parse()
}
