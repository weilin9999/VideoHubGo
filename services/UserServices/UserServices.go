/*
 * @Descripttion: 用户服务层 - User Services
 * @Author: William Wu
 * @Date: 2022/05/28 下午 12:54
 * @LastEditors: William Wu
 * @LastEditTime: 2022/05/28 下午 12:54
 */
package UserServices

import (
	"VideoHubGo/models/UserModel"
	"VideoHubGo/utils/DataBaseUtils"
	"VideoHubGo/utils/EncryptionUtils"
	uuid "github.com/satori/go.uuid"
	"time"
)

var db = DataBaseUtils.GoDB()

/**
 * @Descripttion: 用户登录 - User Login
 * @Author: William Wu
 * @Date: 2022/05/29 下午 06:22
 * @Param: Account (string)
 * @Param: Password (string)
 * @Return: UserModel User
 */
func LoginUser(account string, password string) UserModel.User {
	userSalt := FindUserSalt(account)
	realPwd := EncryptionUtils.ReversePassword(password, userSalt)
	userData := UserModel.User{}
	if err := db.Table("userdata").Where("account = ? and password = ?", account, realPwd).Take(&userData).Error; err != nil {
		return userData
	}
	return userData
}

/**
 * @Descripttion:
 * @Author: William Wu
 * @Date: 2022/05/29 下午 06:30
 * @Param: Account (string)
 * @Param: Password (string)
 * @Return: Result (int)
 */
func RegisterUser(account string, password string, username string) int {
	//先检查是否有人已注册了账号
	isAcAlive, isUnAlive := FindUserAlive(account), FindUserNameAlive(username)
	if isAcAlive {
		return 2
	} else if isUnAlive {
		return 3
	}
	userData := UserModel.UserRegister{}
	pwd, salt := EncryptionUtils.EncPassword(account, password)
	userData.Password = pwd
	userData.Salt = salt
	userData.Username = username
	userData.Account = account
	userData.Isdelete = 0
	userData.Isuploader = 0
	userData.Isadmin = 0
	userData.Avatar = ""
	isDeleteUid := FindIsDeleteUser()
	if isDeleteUid != 0 {
		userData.Uid = isDeleteUid
		userData.Create_Time = time.Now()
		if err := db.Table("userdata").Save(&userData).Error; err != nil {
			return 4
		}
	} else {
		if err := db.Table("userdata").Create(&userData).Error; err != nil {
			return 0
		}
	}
	return 1
}

/**
 * @Descripttion: 更新用户密码 - Update User Password
 * @Author: William Wu
 * @Date: 2022/05/30 下午 06:28
 * @Param: account (string)
 * @Param: oldPassword (string)
 * @Param: newPassword (string)
 * @Return: Result (int)
 */
func UpdatePassword(account string, oldPassword string, newPassword string) int {
	userData := LoginUser(account, oldPassword)
	if userData.Uid == 0 {
		return 0
	}
	encPassword := EncryptionUtils.ReversePassword(newPassword, userData.Salt)
	db.Table("userdata").Where("uid = ?", userData.Uid).Update("password", encPassword)
	return 1
}

/**
 * @Descripttion: 用户头像保存到数据库 - User Avatar Save In DataBase
 * @Author: William Wu
 * @Date: 2022/06/02 上午 10:51
 * @Param: userId (int)
 * @Param: fileTpye (string)
 * @Return: fileName (string)
 */
func UserUploadAvatar(userId int, fileTpye string) string {
	uuidKey := uuid.NewV4()
	nowTime := time.Now().Format("2006-1-2 15:04:05.000")
	uuidValue := uuid.NewV5(uuidKey, nowTime).String()
	db.Table("userdata").Where("uid = ?", userId).Update("avatar", uuidValue+fileTpye)
	return uuidValue + fileTpye
}

/**
 * @Descripttion: 查询账号是否存在 - Query whether the account exists
 * @Author: William Wu
 * @Date: 2022/05/29 下午 06:42
 * @Param: Account (string)
 * @Return: Result (int)
 */
func FindUserAlive(account string) bool {
	tempData := UserModel.User{}
	db.Select("uid").Table("userdata").First(&tempData, "account = ?", account)
	return tempData.Uid != 0
}

/**
 * @Descripttion: 查询账号的盐值 - Query the salt value of the account
 * @Author: William Wu
 * @Date: 2022/05/29 下午 06:42
 * @Param: Account (string)
 * @Return: Salt (string)
 */
func FindUserSalt(account string) string {
	userModel := UserModel.User{}
	if err := db.Table("userdata").Select("salt").Where("account = ?", account).Find(&userModel).Error; err != nil {
		return ""
	}
	return userModel.Salt
}

/**
 * @Descripttion: 查询用户名是否被占用 - Query whether the user name is occupied
 * @Author: William Wu
 * @Date: 2022/05/29 下午 06:42
 * @Param: Username (string)
 * @Return: Result (int)
 */
func FindUserNameAlive(username string) bool {
	tempData := UserModel.User{}
	db.Select("uid").Table("userdata").First(&tempData, "username = ? ", username)
	return tempData.Uid != 0
}

/**
 * @Descripttion: 查找已经被注销删除的用户 - Find Is Deleted User Info
 * @Author: William Wu
 * @Date: 2022/06/09 下午 03:56
 * @Return: uid (int)
 */
func FindIsDeleteUser() int {
	userData := UserModel.User{}
	db.Select("uid").Table("userdata").First(&userData, "isdelete = ?", 1)
	return userData.Uid
}

/**
 * @Descripttion: 根据用户查询是否被锁定 - Whether it is locked according to user query
 * @Author: William Wu
 * @Date: 2022/06/09 下午 04:27
 * @Param: Account (string)
 * @Return: isdelete (int)
 */
func FindIsDeleteUserLogin(account string) int {
	userData := UserModel.User{}
	db.Select("uid").Table("userdata").First(&userData, "account = ?", account)
	return userData.Isdelete
}
