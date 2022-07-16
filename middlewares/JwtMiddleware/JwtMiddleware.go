/*
 * @Descripttion: Jwt Middleware
 * @Author: William Wu
 * @Date: 2022-05-21 16:59:40
 * @LastEditors: William Wu
 * @LastEditTime: 2022-05-21 18:34:51
 */
package JwtMiddleware

import (
	"VideoHubGo/utils/JsonUtils"
	"VideoHubGo/utils/LogUtils"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

/**
 * @Descripttion:  Jwt模型 - Jwt Model
 * @Author: William Wu
 * @Date: 2022/05/23 下午 04:02
 */
type JwtModel struct {
	UserId     int    `json:"uid"`
	Account    string `json:"account"`
	UserName   string `json:"username"`
	IsAdmin    int    `json:"isadmin"`
	IsUploader int    `json:"isuploader"`
	jwt.StandardClaims
}

/**
 * @Descripttion: 生成Token - Create Token
 * @Author: William Wu
 * @Date: 2022/05/21 下午 05:01
 * @Param: 用户UID - User UId
 * @Param: 用户账号 - User Account
 * @Param: 用户昵称 - User Name
 * @Param: 是否为管理员 - Is Admin
 * @Param: 是否为上传者 - Is Uploader
 * @Return: Token JWT
 */
func CreateToken(userId int, account string, userName string, isAdmin int, isUploader int) (string, error) {

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

	key := []byte(config.GetString("jwt.key"))
	issuer := config.GetString("jwt.issuer")
	subject := config.GetString("jwt.subject")
	expireTimeConfig := config.GetInt("jwt.expireTime")

	expireTime := time.Now().Add(time.Duration(expireTimeConfig) * time.Hour) //过期时间 - Expire Time
	nowTime := time.Now()                                                     //当前时间 - Now Time
	claims := JwtModel{
		UserId:     userId,
		Account:    account,
		UserName:   userName,
		IsAdmin:    isAdmin,
		IsUploader: isUploader,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间戳 - Expire Time
			IssuedAt:  nowTime.Unix(),    //当前时间戳 - Now Time
			Issuer:    issuer,            //颁发者签名 - Issuer Signature
			Subject:   subject,           //签名主题   - Signature Theme
		},
	}
	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenStruct.SignedString(key)
}

/**
 * @Descripttion: 验证Token - Check Token
 * @Author: William Wu
 * @Date: 2022/05/23 下午 04:03
 * @Param: token
 * @Return: JwtModel
 */
func CheckToken(token string) (*JwtModel, bool) {
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

	key := []byte(config.GetString("jwt.key"))

	tokenObj, _ := jwt.ParseWithClaims(token, &JwtModel{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if key, _ := tokenObj.Claims.(*JwtModel); tokenObj.Valid {
		return key, true
	} else {
		return nil, false
	}
}

/**
 * @Descripttion: JWT中间件 - JWT Middleware
 * @Author: William Wu
 * @Date: 2022/05/23 下午 04:03
 * @Param: *gin
 * @Return: Middleware Controller
 */
func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从请求头中获取Token - Get Token From Request Header
		tokenStr := c.Request.Header.Get("Authorization")
		//用户不存在 - User Non Existent
		if tokenStr == "" {
			c.JSON(http.StatusOK, JsonUtils.JsonResult(101, "非法请求 - Illegal Request", ""))
			c.Abort() //阻止执行 - Stopping
			return
		}
		//Token格式错误 - Token Format Error
		tokenSlice := strings.SplitN(tokenStr, " ", 2)
		if len(tokenSlice) != 2 && tokenSlice[0] != "Bearer" {
			c.JSON(http.StatusOK, JsonUtils.JsonResult(102, "Token格式错误 - Token Format Error", ""))
			c.Abort() //阻止执行 - Stopping
			return
		}
		//验证Token - Check Token
		tokenStruck, ok := CheckToken(tokenSlice[1])
		if !ok {
			c.JSON(http.StatusOK, JsonUtils.JsonResult(103, "Token错误 - Token Error", ""))
			c.Abort() //阻止执行 - Stopping
			return
		}
		//Token超时 - Token Timeout
		if time.Now().Unix() > tokenStruck.ExpiresAt {
			c.JSON(http.StatusOK, JsonUtils.JsonResult(104, "Token过期 - Token Expired", ""))
			c.Abort() //阻止执行 - Stopping
			return
		}

		c.Set("uid", tokenStruck.UserId)
		c.Set("account", tokenStruck.Account)
		c.Set("username", tokenStruck.UserName)
		c.Set("isadmin", tokenStruck.IsAdmin)
		c.Set("isuploader", tokenStruck.IsUploader)
		c.Next()
	}
}

/**
 * @Descripttion: 后台JWT中间件检测 - Admin JWT Middleware Check
 * @Author: William Wu
 * @Date: 2022/06/30 下午 12:00
 * @Param: *gin
 * @Return: Middleware Controller
 */
func AdminJwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从请求头中获取Token - Get Token From Request Header
		tokenStr := c.Request.Header.Get("Authorization")
		//用户不存在 - User Non Existent
		if tokenStr == "" {
			c.JSON(http.StatusOK, JsonUtils.JsonResult(101, "非法请求 - Illegal Request", ""))
			c.Abort() //阻止执行 - Stopping
			return
		}
		//Token格式错误 - Token Format Error
		tokenSlice := strings.SplitN(tokenStr, " ", 2)
		if len(tokenSlice) != 2 && tokenSlice[0] != "Bearer" {
			c.JSON(http.StatusOK, JsonUtils.JsonResult(102, "Token格式错误 - Token Format Error", ""))
			c.Abort() //阻止执行 - Stopping
			return
		}
		//验证Token - Check Token
		tokenStruck, ok := CheckToken(tokenSlice[1])
		if !ok {
			c.JSON(http.StatusOK, JsonUtils.JsonResult(103, "Token错误 - Token Error", ""))
			c.Abort() //阻止执行 - Stopping
			return
		}
		//Token超时 - Token Timeout
		if time.Now().Unix() > tokenStruck.ExpiresAt {
			c.JSON(http.StatusOK, JsonUtils.JsonResult(104, "Token过期 - Token Expired", ""))
			c.Abort() //阻止执行 - Stopping
			return
		}

		isAdmin := tokenStruck.IsAdmin
		//验证是否为管理员 - Check user is admin
		if isAdmin != 1 {
			c.JSON(http.StatusOK, gin.H{"code": 105, "msg": "非法登录 - Illegal Login"})
			c.Abort() //阻止执行 - Stopping
			return
		}

		c.Set("uid", tokenStruck.UserId)
		c.Set("account", tokenStruck.Account)
		c.Set("username", tokenStruck.UserName)
		c.Set("isadmin", tokenStruck.IsAdmin)
		c.Set("isuploader", tokenStruck.IsUploader)
		c.Next()
	}
}
