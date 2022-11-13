package data

import (
	"fmt"
	"jdb/app"
)

// Tables 表信息
type Tables struct {
	// 表名
	Name string `gorm:"column:NAME" json:"Name"`
	// 字符集
	Collation string `gorm:"column:COLLATION" json:"Collation"`
	// 注释
	Comment string `gorm:"column:COMMENT" json:"Comment"`
}

// Column 字段信息
type Column struct {
	Name string `gorm:"column:NAME" json:"Name"`
	// 默认数据
	Default  string `gorm:"column:DEFAULT" json:"Default"`
	Key      string `gorm:"column:KEY" json:"Key"`
	Nullable string `gorm:"column:NULLABLE" json:"Nullable"`
	// 字段类型 varchar(255)
	Type string `gorm:"column:TYPE" json:"Type"`
	// 注释
	Comment string `gorm:"column:COMMENT" json:"Comment"`
	Extra   string `gorm:"column:EXTRA" json:"Extra"`
}

func (t *Tables) GetByScheme(scheme string) []*Tables {
	var tables []*Tables

	db := app.Instance().GormDB
	if err := db.Table("TABLES").
		Select("TABLE_NAME AS NAME, TABLE_COLLATION AS COLLATION, TABLE_COMMENT AS COMMENT ").
		Where("TABLE_SCHEMA = ?", scheme).
		Order("TABLE_NAME").Find(&tables).Error; err != nil {
		fmt.Println(fmt.Sprintf("GetByScheme:==> find tables for scheme [%s] failed", scheme))
		panic(err)
	}

	return tables
}

func (t *Tables) GetBySchemes(schemes []string) []*Tables {
	var tables []*Tables
	if err := app.Instance().GormDB.Table("TABLES").
		Select("TABLE_NAME AS NAME, TABLE_COLLATION AS COLLATION, TABLE_COMMENT AS COMMENT ").
		Where("TABLE_SCHEMA in (?)", schemes).
		Order("TABLE_NAME").Find(&tables).Error; err != nil {
		fmt.Println(fmt.Sprintf("GetBySchemes:==> find tables for schemes [%s] failed", schemes))
		panic(err)
	}

	return tables
}

// GetBySchemeAndTable 获取库中某个表的字段信息
func (t *Column) GetBySchemeAndTable(scheme, table string) []*Column {
	var columns []*Column

	db := app.Instance().GormDB
	if err := db.Table("COLUMNS").
		Select("COLUMN_NAME AS NAME, COLUMN_DEFAULT AS `DEFAULT`, COLUMN_KEY AS `KEY`, IS_NULLABLE AS NULLABLE, COLUMN_TYPE AS TYPE, COLUMN_COMMENT AS `COMMENT`, EXTRA").
		Where("TABLE_NAME = ? AND TABLE_SCHEMA = ?", table, scheme).
		Order("TABLE_NAME").Find(&columns).Error; err != nil {
		fmt.Println(fmt.Sprintf("GetByScheme:==> find tables for scheme [%s] failed", scheme))
		panic(err)
	}

	return columns
}
