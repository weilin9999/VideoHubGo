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
	Uid         int    `gorm:"primaryKey" json:"uid"`
	Account     string `json:"account"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Salt        string
	Isadmin     int
	Isuploader  int
	Isdelete    int
	Avatar      string
	Create_Time time.Time
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

/**
 * @Descripttion: 用户编辑模型 - User Edit Struct
 * @Author: William Wu
 * @Date: 2022/06/30 下午 10:04
 */
type EditUser struct {
	Uid        int    `gorm:"primaryKey" json:"uid"`
	Account    string `json:"account"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Salt       string
	Avatar     string `json:"avatar"`
	Isadmin    int    `json:"isadmin"`
	Isuploader int    `json:"isuploader"`
	Isdelete   int    `json:"isdelete"`
}

/**
 * @Descripttion: 后台搜索用户模型 - Admin search user struct
 * @Author: William Wu
 * @Date: 2022/07/05 下午 10:57
 */
type UserAdminSearch struct {
	Uid      int    `json:"uid"`
	Account  string `json:"account"`
	Username string `json:"username"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
}

/**
 * @Descripttion: 后台用户模型 - Admin user struct
 * @Author: William Wu
 * @Date: 2022/07/05 下午 10:59
 */
type UserAdminRe struct {
	Uid         int       `json:"uid"`
	Account     string    `json:"account"`
	Username    string    `json:"username"`
	Avatar      string    `json:"avatar"`
	Isadmin     int       `json:"isadmin"`
	Isuploader  int       `json:"isuploader"`
	Isdelete    int       `json:"isdelete"`
	Create_Time time.Time `json:"create_time"`
}
