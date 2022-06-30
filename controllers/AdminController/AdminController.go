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
	"VideoHubGo/models/UserModel"
	"VideoHubGo/services/UserServices"
	"VideoHubGo/utils/JsonUtils"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
 * @Descripttion: 后台登录接口 - Back End Login APi
 * @Author: William Wu
 * @Date: 2022/06/29 下午 11:18
 * @Param: 账号 - Account (string)
 * @Param: 密码 - Password (string)
 * @Return: Token JWT
 */
func AdminLogin(ctx *gin.Context) {
	userData := UserModel.User{}
	err := ctx.BindJSON(&userData)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	if userData.Account == "" || userData.Password == "" {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "请输入正确的账号密码 - Please enter the correct account password\n\n", ""))
		return
	}
	isAlive := UserServices.FindUserAlive(userData.Account)
	if isAlive {
		userData = UserServices.LoginUser(userData.Account, userData.Password)
	} else {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(202, "该用户不存在 - This account is not alive", ""))
		return
	}

	isDelete := UserServices.FindIsDeleteUserLogin(userData.Account)
	if isDelete == 1 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(206, "该用户已被锁定 - The user is locked", ""))
		return
	}

	if userData.Uid == 0 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(204, "密码错误 - Password Error", ""))
		return
	}

	token, err := JwtMiddleware.CreateToken(userData.Uid, userData.Username, userData.Isadmin, userData.Isuploader)

	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(205, "Token生成失败 - Token Create Error", ""))
		return
	}

	rData := map[string]interface{}{"token": token, "username": userData.Username, "avatar": userData.Avatar, "upload": userData.Isuploader, "account": userData.Account}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", rData))
}
