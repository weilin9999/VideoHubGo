/*
 * @Descripttion: 用户模型 - User Model
 * @Author: William Wu
 * @Date: 2022-05-21 22:03:49
 * @LastEditors: William Wu
 * @LastEditTime: 2022-05-21 22:03:55
 */
package UserModel

import (
	"time"
)

/**
 * @Descripttion: 用户模型 - User Model
 * @Author: William Wu
 * @Date: 2022/05/23 下午 03:56
 */
type User struct {
	Uid         int    `gorm:"primaryKey" json:"uid"`
	Account     string `json:"account"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Salt        string
	Avatar      string
	Isadmin     int
	Isuploader  int
	Isdelete    int
	Create_Time time.Time
	Update_Time time.Time
}

/**
 * @Descripttion: 用户注册模型 - User Register Model
 * @Author: William Wu
 * @Date: 2022/05/23 下午 03:56
 */
type UserRegister struct {
	Uid      int    `gorm:"primaryKey" json:"uid"`
	Account  string `json:"account"`
	Username string `json:"username"`
	Password string `json:"password"`
	Salt     string
}

/**
 * @Descripttion: 用户更新密码模型 - User Update Password Model
 * @Author: William Wu
 * @Date: 2022/05/23 下午 03:56
 */
type UserUpdatePassword struct {
	Account     string `json:"account"`
	Password    string `json:"password"`
	RePassword  string `json:"repassword"`
	NewPassword string `json:"newpassword"`
}
