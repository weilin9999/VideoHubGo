/*
 * @Descripttion: 管理员控制层 - Admin Controller
 * @Author: William Wu
 * @Date: 2022-05-21 11:51:55
 * @LastEditors: William Wu
 * @LastEditTime: 2022-05-21 18:35:03
 */
package AdminController

import (
	"VideoHubGo/middlewares/JwtMiddleware"

	"github.com/gin-gonic/gin"
)

func AdminLogin(ctx *gin.Context) {
	username := JwtMiddleware.GetTokenUsername(ctx)
	dta := map[string]interface{}{
		"tag":      "<br>",
		"username": username,
	}
	ctx.AsciiJSON(200, dta)
}
