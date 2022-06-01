/*
 * @Descripttion: 日志工具 - Log Utils
 * @Author: William Wu
 * @Date: 2022/05/27 下午 11:13
 * @LastEditors: William Wu
 * @LastEditTime: 2022/05/27 下午 11:13
 */
package LogUtils

import (
	"log"
	"os"
)

func init() {
	fileName, err := os.OpenFile("./logger.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Logger(err.Error())
	}

	log.SetOutput(fileName)
}

func Logger(anything string) {
	log.Println(anything)
}
