package core

import (
	"fmt"
	"multipart-upload/global"

	"github.com/spf13/viper"
)

// InitConfig /* 初始化配置文件
func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("viper.ReadInConfig err=%s\n", err))
	}

	if err := viper.Unmarshal(&global.Config); err != nil {
		panic(fmt.Errorf("viper.Unmarshal err=%s\n", err))
	}
}
