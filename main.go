/*
 * @Descripttion: Main Function
 * @Author: William Wu
 * @Date: 2022-05-20 18:15:49
 * @LastEditors: William Wu
 * @LastEditTime: 2022-05-21 19:15:35
 */
package main

import (
	"VideoHubGo/utils/LogUtils"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	router "VideoHubGo/routers"
)

/**
 * @Descripttion: 主函数 - Main Function
 * @Author: William Wu
 * @Date: 2022/05/23 下午 03:59
 */
func main() {
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
	port := config.GetString("http.port")

	//Gin服务启动 - Running Gin Service
	r := gin.Default()
	r = router.Router(r)
	r.Run(":" + port)
}
