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
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"os"
	"time"
)

var RedisClient *redis.Client

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
	db := config.GetInt("redis.db")
	maxActive := config.GetInt("redis.maxActive")
	//拼接连接地址 - Splicing Connection Address
	server := fmt.Sprintf("%s:%s", host, port)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     server,
		Password: password,
		DB:       db,
		PoolSize: maxActive,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		LogUtils.Logger("[Redis启动失败-Redis start fail] 错误-Error：" + err.Error())
	}

}
