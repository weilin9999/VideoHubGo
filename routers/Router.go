/*
 * @Descripttion: Router Manager
 * @Author: William Wu
 * @Date: 2022-05-20 18:50:21
 * @LastEditors: William Wu
 * @LastEditTime: 2022-05-21 17:54:58
 */
package router

import (
	"VideoHubGo/controllers/AdminController"
	"VideoHubGo/controllers/UserController"
	"VideoHubGo/controllers/VideoController"
	"VideoHubGo/middlewares/JwtMiddleware"

	"github.com/gin-gonic/gin"
)

/**
 * @Descripttion: 路由管理 - Router Manager
 * @Author: William Wu
 * @Date: 2022/05/23 下午 03:59
 * @Param: router(gin.Engine)
 * @Return: router
 */
func Router(router *gin.Engine) *gin.Engine {

	//用户路由 - User Router
	routerList1 := router.Group("/user")
	{
		routerList1.POST("/login", UserController.UserLogin)       //用户登录控制器 - User Login Controller
		routerList1.POST("/register", UserController.UserRegister) //用户注册控制器 - User Register Controller
	}

	//用户信息路由 - User Information Route
	routerList2 := router.Group("/info") //.Use(JwtMiddleware.JwtMiddleware()) // JWT中间件 - JWT Middleware
	{
		routerList2.POST("/updatepassword", UserController.UserUpdatePassword) //用户修改密码控制器 - User Update Password Controller
	}

	//后台路由 - Admin Route
	routerList3 := router.Group("/admin").Use(JwtMiddleware.JwtMiddleware()) // JWT中间件 - JWT Middleware
	{
		routerList3.POST("/login", AdminController.AdminLogin) //用户登录控制器 - User Login Controller
	}

	//视频路由 - Video Route
	routerList4 := router.Group("/video") //.Use(JwtMiddleware.JwtMiddleware()) // JWT中间件 - JWT Middleware
	{
		routerList4.POST("/list", VideoController.GetVideoList) //视频控制器 - Video Controller
	}

	return router
}
