package base

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

// MySQLConf
type MySQLConf struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DB       string `yaml:"db"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// 连接数据库
func NewMySQLConnection(username, password, protocol, host, port, dbname string) (*gorm.DB, error) {

	// username:password@protocol(address)/dbname?param=value
	dbName := fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, protocol,
		fmt.Sprintf("%s:%s", host, port), dbname)

	database, err := gorm.Open("mysql", dbName)
	if err != nil {
		return nil, err
	}

	underlyingDB := database.DB()
	if underlyingDB == nil {
		return nil, errors.New("underlying DB is nil")
	}

	underlyingDB.SetMaxOpenConns(20)
	underlyingDB.SetMaxIdleConns(20)
	underlyingDB.SetConnMaxLifetime(time.Duration(600 * int64(time.Millisecond)))

	return database, nil
}

// 创建数据库，存在则忽略
func checkCreateDatabase(username string, password string, protocol string, address string, dbname string) error {
	sqlHostInfo := fmt.Sprintf("%s:%s@%s(%s)/", username, password, protocol, address)
	db, err := sql.Open("mysql", sqlHostInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("create database if not exists " + dbname)
	if err != nil {
		return err
	}
	return nil
}
