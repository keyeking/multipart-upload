package main

import (
	"fmt"
	"multipart-upload/core"
	"multipart-upload/global"
)

func main() {
	// 初始化配置
	core.InitConfig()
	core.InitFlag()
	core.InitGorm()
	if global.Config.Gorm.IsMigrate {
		return
	}
	core.InitRouter()
	err := global.Router.Run(fmt.Sprintf("%s:%s", global.Config.System.Host, global.Config.System.Port))
	if err != nil {
		fmt.Println(err)
		return
	}
}
