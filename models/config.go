package models

import "fmt"

type Config struct {
	System System `yaml:"system"`
	File   File   `yaml:"file"`
	Gorm   Gorm   `yaml:"gorm"`
}

type System struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}
type File struct {
	Path string `yaml:"path"`
}

type Gorm struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Db        string `yaml:"db"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	IsMigrate bool   //是否进行数据库迁移
}

func (g *Gorm) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", g.User, g.Password, g.Host, g.Port, g.Db)
}
