/*
 * @Descripttion: 数据库工具 - DataBase Utils
 * @Author: William Wu
 * @Date: 2022-05-21 11:37:13
 * @LastEditors: William Wu
 * @LastEditTime: 2022-05-21 13:18:21
 */
package DataBaseUtils

import (
	"VideoHubGo/utils/LogUtils"
	"fmt"
	"os"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _db *gorm.DB

/**
 * @Descripttion: 数据库连接 - DataBase Connection
 * @Author: William Wu
 * @Date: 2022/05/23 下午 04:06
 */
func init() {
	//读取配置文件 - Read The Configuration File
	path, err := os.Getwd()
	if err != nil {
		LogUtils.Logger(err.Error())
	}
	config := viper.New()
	config.AddConfigPath(path + "/configs")
	config.SetConfigName("config")
	config.SetConfigType("yaml")

	//尝试进行配置读取 - Try Reading Configuration
	if err := config.ReadInConfig(); err != nil {
		LogUtils.Logger(err.Error())
	}

	host := config.GetString("database.host")
	port := config.GetString("database.port")
	databasename := config.GetString("database.databasename")
	username := config.GetString("database.username")
	password := config.GetString("database.password")

	//拼接MySQL连接地址 - Splicing Mysql Connection Address
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, databasename)
	// fmt.Println(url)
	var errdb error

	_db, errdb = gorm.Open(mysql.Open(url), &gorm.Config{})
	if errdb != nil {
		LogUtils.Logger(errdb.Error())
	}

	// sqlDB, _ := _db.DB()

	// sqlDB.SetMaxIdleConns(100) //设置最大连接数 - Set Max SQL Connection
	// sqlDB.SetMaxIdleConns(20)  //设置最大空闲连接数 - Set Max Free Connection

}

/**
 * @Descripttion: 数据库连接对象 - DataBase Connection Object
 * @Author: William Wu
 * @Date: 2022/05/23 下午 04:06
 */
func GoDB() *gorm.DB {
	return _db
}
