/*
 * @Descripttion: 无效路由拦截器 - No Route Middlewares
 * @Author: William Wu
 * @Date: 2022/06/02 下午 04:56
 * @LastEditors: William Wu
 * @LastEditTime: 2022/06/02 下午 04:56
 */
package NoRouteMiddleware

import (
	"VideoHubGo/utils/JsonUtils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NoRouteHttp(ctx *gin.Context) {
	//返回信息 - Return Message
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(400, "无效访问 - Invalid access", ""))
	return
}

func NoMethodHttp(ctx *gin.Context) {
	//返回信息 - Return Message
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(400, "不允许的方法 - Method Not Allowed", ""))
	return
}
