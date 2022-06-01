/*
 * @Descripttion: Redis工具 - Redis Utils
 * @Author: William Wu
 * @Date: 2022/05/27 下午 05:11
 * @LastEditors: William Wu
 * @LastEditTime: 2022/05/27 下午 05:11
 */
package RedisUtils

import (
	"VideoHubGo/utils/LogUtils"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"os"
	"time"
)

var RedisPool *redis.Pool

/**
 * @Descripttion: Redis连接 - Redis Connection
 * @Author: William Wu
 * @Date: 2022/05/27 下午 10:55
 * @Return: Redis
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

	host := config.GetString("redis.host")
	port := config.GetString("redis.port")
	password := config.GetString("redis.password")
	maxIdle := config.GetInt("redis.maxIdle")
	maxActive := config.GetInt("redis.maxActive")
	idleTimeout := config.GetInt("redis.idleTimeout")
	option := redis.DialPassword(password)
	//拼接连接地址 - Splicing Connection Address
	server := fmt.Sprintf("%s:%s", host, port)

	RedisPool = &redis.Pool{
		MaxIdle:     maxIdle,                    //最初的连接数量 - Number Of Initial Connections
		MaxActive:   maxActive,                  //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配 - Max Connections
		IdleTimeout: time.Duration(idleTimeout), //连接关闭时间 300秒 （300秒不使用自动关闭） - Connections Timeout
		Dial: func() (redis.Conn, error) { //要连接的Redis数据库 - Need To Connection Redis DataBase
			conn, err := redis.Dial("tcp", server, option)
			if err != nil {
				conn.Close()
			}
			return conn, err
		},
	}
}
