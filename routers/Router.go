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
	"VideoHubGo/controllers/ClassController"
	"VideoHubGo/controllers/UserController"
	"VideoHubGo/controllers/VideoController"
	"VideoHubGo/middlewares/JwtMiddleware"
	"VideoHubGo/middlewares/NoRouteMiddlewares"
	"VideoHubGo/utils/LogUtils"
	"VideoHubGo/utils/UploadUtils"
	"github.com/spf13/viper"
	"os"

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

	router.NoRoute(NoRouteMiddlewares.NoRouteHttp)   //设置网页错误返回信息 - Setting Http Error Return Message
	router.NoMethod(NoRouteMiddlewares.NoMethodHttp) //设置网页错误返回信息 - Setting Http Error Return Message

	maxSize := config.GetInt("files.maxSize") //设置最大文件上传大小 - Set Max Upload File Size

	router.MaxMultipartMemory = int64(maxSize << 20)

	//用户路由 - User Router
	routerList1 := router.Group("/user")
	{
		routerList1.POST("/login", UserController.UserLogin)       //用户登录控制器 - User Login Controller
		routerList1.POST("/register", UserController.UserRegister) //用户注册控制器 - User Register Controller
	}

	//用户信息路由 - User Information Route
	routerList2 := router.Group("/userinfo").Use(JwtMiddleware.JwtMiddleware()) // JWT中间件 - JWT Middleware
	{
		routerList2.POST("/update/password", UserController.UserUpdatePassword) //用户修改密码控制器 - User Update Password Controller
		routerList2.POST("/upload/avatar", UserController.UploadAvatar)         //用户头像上传控制器 - User Upload Avatar Controller
	}

	//后台路由 - Admin Route
	routerList3 := router.Group("/admin").Use(JwtMiddleware.JwtMiddleware()) // JWT中间件 - JWT Middleware
	{
		routerList3.POST("/login", AdminController.AdminLogin) //用户登录控制器 - User Login Controller
	}

	//视频分类路由 - Video Class Route
	routerList4 := router.Group("/class") //.Use(JwtMiddleware.JwtMiddleware()) // JWT中间件 - JWT Middleware
	{
		routerList4.POST("/list", ClassController.GetClassList) //视频控制器 - Video Controller
	}

	//视频路由 - Video Route
	routerList5 := router.Group("/video") //.Use(JwtMiddleware.JwtMiddleware()) // JWT中间件 - JWT Middleware
	{
		routerList5.POST("/list", VideoController.GetVideoList)            //总视频控制器 - Center Video Controller
		routerList5.POST("/class/list", VideoController.GetVideoClassList) //视频类型控制器 - Video Class Controller
	}

	//文件映射 - Map File
	routerList6 := router.Group("/file")
	{
		routerList6.Static("/avatar", UploadUtils.GetUploadFilePath("user.userAvatar")) //映射头像文件夹 - Map avatar folder
	}

	return router
}
