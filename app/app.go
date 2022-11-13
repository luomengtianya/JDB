package app

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"jdb/base"
)

var instance *App

// App 配置信息
type App struct {
	MySQL    base.MySQLConf `yaml:"mysql"`
	JGetConf *JGetConf      `yaml:"jget"`

	GormDB       *gorm.DB
	GormDBForAll *gorm.DB
}

// JGetConf 配置信息
type JGetConf struct {
	// 排除的数据
	Exclude []string `yaml:"exclude"`
	// 需要处理的数据
	Include []string `yaml:"include"`
	// 需要处理的库数据
	Scheme []string `yaml:"scheme"`
	// 文件输出路径
	Out string `yaml:"out"`
}

func init() {
	instance = &App{
		JGetConf: &JGetConf{},
	}
}

// Instance get App instance
func Instance() *App {
	return instance
}

// InitGormDB 连接数据库
func (c *App) InitGormDB() error {

	db, err := base.NewMySQLConnection(c.MySQL.User, c.MySQL.Password,
		"tcp", c.MySQL.Host, c.MySQL.Port, c.MySQL.DB)
	if err != nil {
		return fmt.Errorf("new mysql connection error, err: %v", err)
	}

	db.LogMode(true)
	db.Debug()

	c.GormDB = db
	return nil
}
