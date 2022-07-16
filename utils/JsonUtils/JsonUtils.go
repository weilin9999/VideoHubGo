/*
 * @Descripttion: Json工具 -Json Utils
 * @Author: William Wu
 * @Date: 2022/05/29 下午 05:46
 * @LastEditors: William Wu
 * @LastEditTime: 2022/05/29 下午 05:46
 */
package JsonUtils

import (
	"github.com/gin-gonic/gin"
)

func JsonResult(code int, msg string, data any) gin.H {
	h := gin.H{"code": code, "msg": msg, "data": data}
	return h
}
