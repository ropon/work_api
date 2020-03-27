package dblayer

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"workApi/config"
)

var DB *gorm.DB

// InitMySQL 初始化MySQL连接
func InitMySQL() (err error) {
	// Creating a connection to the database
	DB, err = gorm.Open("mysql", config.GetDSN())
	if err != nil {
		return
	}
	err = DB.DB().Ping()
	return err
}
