/*
 * @Descripttion: Router Manager
 * @Author: William Wu
 * @Date: 2022/06/01 下午 11:52
 * @LastEditors: William Wu
 * @LastEditTime: 2022/06/01 下午 11:52
 */
package UploadUtils

import (
	"VideoHubGo/utils/LogUtils"
	"github.com/spf13/viper"
	"os"
)

/**
 * @Descripttion: 获取配置文件 - Get Configuration Path
 * @Author: William Wu
 * @Date: 2022/06/01 下午 11:55
 * @Param: where (string)
 * @Return: Path (string)
 */
func GetFilePath(where string) string {
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

	filePath := config.GetString("files." + where)

	return filePath
}
