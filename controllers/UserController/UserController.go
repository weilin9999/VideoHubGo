/*
 * @Descripttion: 用户控制层 - User Controller
 * @Author: William Wu
 * @Date: 2022-05-20 22:04:49
 * @LastEditors: William Wu
 * @LastEditTime: 2022-05-21 22:04:35
 */

package UserController

import (
	"VideoHubGo/middlewares/JwtMiddleware"
	"VideoHubGo/models/UserModel"
	"VideoHubGo/services/UserServices"
	"VideoHubGo/utils/DataBaseUtils"
	"VideoHubGo/utils/JsonUtils"
	"VideoHubGo/utils/UploadUtils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
)

var db = DataBaseUtils.GoDB()

/**
 * @Descripttion: 登录控制器 - Login Controller
 * @Author: William Wu
 * @Date: 2022/05/23 下午 04:06
 * @Param: 账号 - Account (string)
 * @Param: 密码 - Password (string)
 * @Return: Token JWT
 */
func UserLogin(ctx *gin.Context) {

	userData := UserModel.User{}
	ctx.Bind(&userData)
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

	if userData.Uid == 0 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(204, "密码错误 - Password Error", ""))
		return
	}

	token, err := JwtMiddleware.CreateToken(userData.Uid, userData.Username, userData.Isadmin, userData.Isuploader)

	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(205, "Token生成失败 - Token Create Error", ""))
		return
	}

	rData := map[string]interface{}{"token": token}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", rData))
}

/**
 * @Descripttion: 注册控制器 - Register Controller
 * @Author: William Wu
 * @Date: 2022/05/28 上午 09:14
 * @Param: UserModel UserRegister
 * @Return: Json
 */
func UserRegister(ctx *gin.Context) {
	userEntity := UserModel.UserRegister{}
	ctx.Bind(&userEntity)
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{4,16}$", userEntity.Account); !ok {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(205, "账号不符合要求! - Account is not allow", ""))
		return
	}
	const nicknamePattern = `^[a-z0-9A-Z\p{Han}]+(_[a-z0-9A-Z\p{Han}]+)*$`
	var nicknameRegexp = regexp.MustCompile(nicknamePattern)
	if len(userEntity.Username) < 4 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(206, "用户名不符合要求! - Username is not allow", ""))
		return
	}
	if !nicknameRegexp.MatchString(userEntity.Username) {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(207, "用户名不符合要求! - Username is not allow", ""))
		return
	}
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{4,16}$", userEntity.Password); !ok {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(208, "密码不符合要求! - Password is not allow", ""))
		return
	}
	res := UserServices.RegisterUser(userEntity.Account, userEntity.Password, userEntity.Username)
	if res == 1 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "注册成功! - Register Success", ""))
	} else if res == 2 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(202, "注册账号已存在 - Register Account Already Exists ", ""))
	} else if res == 3 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(203, "注册用户名已存在 - Register Username Already Exists ", ""))
	} else {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(204, "注册时产生异常 - Exception Occurred During Registration", ""))
	}
}

/**
 * @Descripttion: 用户更改密码 - Update User Password
 * @Author: William Wu
 * @Date: 2022/06/01 下午 07:38
 * @Param: account (string)
 * @Param: oldPassword (string)
 * @Param: rePassword (string)
 * @Param: newPassword (string)
 * @Return: Json
 */
func UserUpdatePassword(ctx *gin.Context) {
	userData := UserModel.UserUpdatePassword{}
	ctx.Bind(&userData)
	if userData.Password != userData.RePassword {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "两次密码不一致! - Second is Inconsistent", ""))
		return
	}
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{4,16}$", userData.NewPassword); !ok {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(202, "密码不符合要求! - Password is not allow", ""))
		return
	}
	res := UserServices.UpdatePassword(userData.Account, userData.Password, userData.NewPassword)
	if res == 1 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "密码更改成功! - Update Password Success", ""))
	} else {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(203, "密码更改时产生异常 - Exception Occurred During Update Password", ""))
	}
}

/**
 * @Descripttion: 用户上传头像 - User Upload Avatar
 * @Author: William Wu
 * @Date: 2022/06/01 下午 07:53
 * @Param: File
 * @Return: Json
 */
func UploadAvatar(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "未上传文件 - Is not upload file", ""))
		return
	} else {
		if file.Size > 2097152 { //小于2M文件 - Less Than 2M File Size
			ctx.JSON(http.StatusOK, JsonUtils.JsonResult(202, "文件大于2M - Files larger than 2M Size", ""))
			return
		}
		fileType := path.Ext(file.Filename)
		if fileType != ".png" && fileType != ".jpg" && fileType != ".jpeg" && fileType != ".gif" {
			fmt.Println(fileType)
			ctx.JSON(http.StatusOK, JsonUtils.JsonResult(203, "文件类型不符合要求 - Document type does not meet requirements", ""))
			return
		}
		userID := JwtMiddleware.GetTokenUID(ctx)
		fileName := UserServices.UserUploadAvatar(userID, fileType)
		filePath := filepath.Join(UploadUtils.GetUploadFilePath("user.userAvatar"), fileName)
		ctx.SaveUploadedFile(file, filePath)
		rData := map[string]interface{}{"avatar": fileName}
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "上传头像成功 - Upload Avatar Success", rData))
	}
}
