/*
 * @Descripttion: 后台服务层 - Admin Services
 * @Author: William Wu
 * @Date: 2022/06/30 上午 11:40
 * @LastEditors: William Wu
 * @LastEditTime: 2022/06/30 上午 11:40
 */
package AdminServices

import (
	"VideoHubGo/models/AdminModel"
	"VideoHubGo/models/ClassModel"
	"VideoHubGo/models/PhotoModel"
	"VideoHubGo/models/UserModel"
	"VideoHubGo/models/VideoModel"
	"VideoHubGo/services/UserServices"
	"VideoHubGo/utils/DataBaseUtils"
	"VideoHubGo/utils/EncryptionUtils"
	"VideoHubGo/utils/LogUtils"
	"time"
)

var db = DataBaseUtils.GoDB()

/**
 * @Descripttion: 管理员登录 - Admin Login
 * @Author: William Wu
 * @Date: 2022/06/30 上午 11:58
 * @Param: Account (string)
 * @Param: Password (string)
 * @Return: UserModel User
 */
func LoginAdmin(account string, password string) UserModel.User {
	userSalt := UserServices.FindUserSalt(account)
	realPwd := EncryptionUtils.ReversePassword(password, userSalt)
	userData := UserModel.User{}
	if err := db.Table("userdata").Where("account = ? AND password = ? AND isadmin = ? AND isdelete = ?", account, realPwd, 1, 0).Take(&userData).Error; err != nil {
		return userData
	}
	return userData
}

/**
 * @Descripttion: 后台搜索视频 - Admin search video
 * @Author: William Wu
 * @Date: 2022/07/05 下午 01:00
 * @Param: vid (string)
 * @Param: detail (string)
 * @Param: cid (string)
 * @Param: size (int)
 * @Param: page (int)
 * @Return: videoData (videoModel.VideoAdminRe), count (int)
 */
func GetVideoList(vid string, detail string, cid int, size int, offset int) ([]VideoModel.VideoAdminRe, int) {
	var videoData []VideoModel.VideoAdminRe
	var count int64
	querySQL := db.Table("videodata")
	if cid != 0 {
		querySQL.Where("cid = ?", cid)
	}
	if vid != "" {
		querySQL.Where("vid LIKE ?", "%"+vid+"%")
	}
	if detail != "" {
		querySQL.Where("detail LIKE ?", "%"+detail+"%")
	}
	querySQL.Where("isdelete = ?", 0).Order("vid DESC").Limit(size).Offset(offset).Count(&count).Find(&videoData)
	return videoData, int(count)
}

/**
 * @Descripttion: 搜索未分类的视频 - Seach no cid video data
 * @Author: William Wu
 * @Date: 2022/07/08 下午 06:04
 * @Param: size (int)
 * @Param: page (int)
 * @Return: videoData (videoModel.VideoAdminRe), count (int)
 */
func GetNoCidVideoList(size int, offset int) ([]VideoModel.VideoAdminRe, int) {
	var videoData []VideoModel.VideoAdminRe
	var count int64
	db.Table("videodata").Where("isdelete = ? AND cid = ?", 0, 0).Order("vid DESC").Limit(size).Offset(offset).Count(&count).Find(&videoData)
	return videoData, int(count)
}

/**
 * @Descripttion: 视频信息修改 - Update Video Informatica
 * @Author: William Wu
 * @Date: 2022/06/30 下午 09:26
 * @Param: VideoEdit (struct)
 * @Return: Result (int)
 */
func EditVideoInformation(VideoData VideoModel.VideoEdit) int {
	if err := db.Table("videodata").Where("vid = ?", VideoData.Vid).Updates(&VideoData).Error; err != nil {
		return 0
	}
	return 1
}

/**
 * @Descripttion: 编辑视频分类信息 - Edit video cid information
 * @Author: William Wu
 * @Date: 2022/07/05 下午 08:47
 * @Param: cid (int)
 * @Param: vid (int)
 * @Return: Result (int)
 */
func EditVideoCidInformation(cid int, vid int) int {
	if err := db.Table("videodata").Where("vid = ?", vid).Update("cid", cid).Error; err != nil {
		return 0
	}
	return 1
}

/**
 * @Descripttion: 软删除视频 - Delete video
 * @Author: William Wu
 * @Date: 2022/06/30 下午 09:51
 * @Param: vid (int)
 * @Return: Result (int)
 */
func DeleteVideoInformation(vid int) (int, string) {
	videoData := VideoModel.VideoEdit{}
	if err := db.Table("videodata").Where("vid = ?", vid).Update("isdelete", "1").Error; err != nil {
		return 0, ""
	}
	if err := db.Select("detail").Table("videodata").Where("vid = ?", vid).First(&videoData).Error; err != nil {
		return 0, ""
	}
	return 1, videoData.Detail
}

/**
 * @Descripttion: 后台获取用户数据 - Admin get user data
 * @Author: William Wu
 * @Date: 2022/07/05 下午 11:04
 * @Param: uid (int)
 * @Param: account (string)
 * @Param: username (string)
 * @Param: size (int)
 * @Param: page (int)
 * @Return: classData (ClassModel.ClassAdminRe), count (int)
 */
func GetUserList(uid int, account string, username string, size int, offset int) ([]UserModel.UserAdminRe, int) {
	var userData []UserModel.UserAdminRe
	var count int64
	querySQL := db.Table("userdata")
	if uid != 0 {
		querySQL.Where("uid = ?", uid)
	}
	if account != "" {
		querySQL.Where("account LIKE ?", "%"+account+"%")
	}
	if username != "" {
		querySQL.Where("username LIKE ?", "%"+username+"%")
	}
	querySQL.Where("isdelete = ?", 0).Limit(size).Offset(offset).Count(&count).Find(&userData)
	return userData, int(count)
}

/**
 * @Descripttion: 用户信息修改 - Update User Information
 * @Author: William Wu
 * @Date: 2022/06/30 下午 10:03
 * @Param: EditUser (struct)
 * @Return: Result (int)
 */
func EditUserInformation(userData UserModel.EditUser) int {
	if userData.Password != "undefined" {
		password, err := UpdatePassword(userData.Uid, userData.Password)
		if err == 0 {
			return 2
		}
		userData.Password = password
	} else {
		userData.Password = ""
	}
	if err := db.Table("userdata").Where("uid = ?", userData.Uid).Updates(&userData).Error; err != nil {
		return 0
	}
	return 1
}

/**
 * @Descripttion: 提升/下降用户权限 - Promote or Decline User authority
 * @Author: William Wu
 * @Date: 2022/07/05 下午 11:23
 * @Param: uid (int)
 * @Param: isadmin (int)
 * @Return: Result (int)
 */
func AuthorityUserIsadmin(uid int, isadmin int) int {
	if err := db.Table("userdata").Where("uid = ?", uid).Update("isadmin", isadmin).Error; err != nil {
		return 0
	}
	return 1
}

/**
 * @Descripttion: 提升/下降用户权限 - Promote or Decline User authority
 * @Author: William Wu
 * @Date: 2022/07/05 下午 11:24
 * @Param: uid (int)
 * @Param: isuploader (int)
 * @Return: Result (int)
 */
func AuthorityUserIsuploader(uid int, isuploader int) int {
	if err := db.Table("userdata").Where("uid = ?", uid).Update("isuploader", isuploader).Error; err != nil {
		return 0
	}
	return 1
}

/**
 * @Descripttion: 软删除用户 - Delete User Information
 * @Author: William Wu
 * @Date: 2022/06/30 下午 10:39
 * @Param: uid (int)
 * @Return: Result (int)
 */
func DeleteUserInformation(uid int) int {
	if err := db.Table("userdata").Where("uid = ?", uid).Update("isdelete", "1").Error; err != nil {
		return 0
	}
	return 1
}

/**
 * @Descripttion: 创建新视频分类 - Created new video class
 * @Author: William Wu
 * @Date: 2022/07/09 上午 11:11
 * @Param: classname (string)
 * @Return: Result (int)
 */
func CreatedClassInformarion(classname string) int {
	classData := ClassModel.Class{}
	classData.Create_Time = time.Now()
	classData.Update_Time = time.Now()
	classData.Isdelete = 0
	classData.Classname = classname
	isDeleteCid := IsDeleteClassInformation()
	if isDeleteCid != 0 {
		classData.Cid = isDeleteCid
		if err := db.Table("vclass").Save(&classData).Error; err != nil {
			LogUtils.Logger("[数据库错误 SQL Error]在操作插入新视频分类更新时产生异常：" + err.Error())
			return 2
		}
	} else {
		if err := db.Table("vclass").Create(&classData).Error; err != nil {
			LogUtils.Logger("[数据库错误 SQL Error]在操作插入新视频分类创建时产生异常：" + err.Error())
			return 3
		}
	}
	return 1
}

/**
 * @Descripttion: 查询已被删除的视频分类 - Query isdelete video class
 * @Author: William Wu
 * @Date: 2022/07/09 上午 11:12
 * @Return: cid (int)
 */
func IsDeleteClassInformation() int {
	classData := ClassModel.Class{}
	db.Select("cid").Table("vclass").First(&classData, "isdelete = ?", 1)
	return classData.Cid
}

/**
 * @Descripttion: 后台搜索获取视频分类数据 - Admin get video class data
 * @Author: William Wu
 * @Date: 2022/07/05 下午 09:21
 * @Param: cid (int)
 * @Param: classname (string)
 * @Param: size (int)
 * @Param: page (int)
 * @Return: classData (ClassModel.ClassAdminRe), count (int)
 */
func GetClassList(cid int, classname string, size int, offset int) ([]ClassModel.ClassAdminRe, int) {
	var classData []ClassModel.ClassAdminRe
	var count int64
	querySQL := db.Table("vclass")
	if cid != 0 {
		querySQL.Where("cid = ?", cid)
	}
	if classname != "" {
		querySQL.Where("classname LIKE ?", "%"+classname+"%")
	}
	querySQL.Where("isdelete = ?", 0).Order("cid DESC").Limit(size).Offset(offset).Count(&count).Find(&classData)
	return classData, int(count)
}

/**
 * @Descripttion: 修改视频分类信息 - Update Video Class Information
 * @Author: William Wu
 * @Date: 2022/06/30 下午 11:25
 * @Param: ClassRe (struct)
 * @Return: Result (int)
 */
func EditClassInformation(classData ClassModel.ClassRe) int {
	if err := db.Table("vclass").Where("cid = ?", classData.Cid).Updates(&classData).Error; err != nil {
		return 0
	}
	return 1
}

/**
 * @Descripttion: 软删除视频分类信息 - Delete Video Class Information
 * @Author: William Wu
 * @Date: 2022/06/30 下午 11:26
 * @Param: cid (int)
 * @Return: Result (int)
 */
func DeleteClassInformation(cid int) int {
	if err := db.Table("vclass").Where("cid = ?", cid).Update("isdelete", "1").Error; err != nil {
		return 0
	}
	return 1
}

/**
 * @Descripttion: 后台更改用户密码 - Admin edit user password
 * @Author: William Wu
 * @Date: 2022/07/06 上午 10:50
 * @Param: uid (int)
 * @Param: password (string)
 * @Return: Reslut (int)
 */
func UpdatePassword(uid int, newPassword string) (string, int) {
	userData := UserModel.User{}
	db.Table("userdata").Select("salt").Where("uid = ?", uid).Find(&userData)
	encPassword := EncryptionUtils.ReversePassword(newPassword, userData.Salt)
	if err := db.Table("userdata").Where("uid = ?", uid).Update("password", encPassword).Error; err != nil {
		return "", 0
	}
	return userData.Password, 1
}

/**
 * @Descripttion: 获取后台二次编码状态 - Get back-end recode status
 * @Author: William Wu
 * @Date: 2022/07/06 下午 06:59
 * @Return: ReCode (int)
 */
func GetReCodeStatus() int {
	reCode := AdminModel.AdminDashboard{}
	db.Table("dashboard").Where("id = ?", 1).Find(&reCode)
	return reCode.Re_Code
}

/**
 * @Descripttion: 更改后台二次编码状态 - Change back-end recode status
 * @Author: William Wu
 * @Date: 2022/07/06 下午 07:00
 * @Return: Result (int)
 */
func ChangeReCodeStatus(status int) int {
	if err := db.Table("dashboard").Where("id = ?", 1).Update("re_code", status).Error; err != nil {
		return 0
	}
	return 1
}

/**
 * @Descripttion: 删除指定图片 - Delete photo
 * @Author: William Wu
 * @Date: 2022/07/09 上午 12:04
 * @Param: pid (int)
 * @Return: Result (int)
 */
func DeletePhotoInformation(pid int) (int, string) {
	photoData := PhotoModel.PhotoModelRe{}
	if err := db.Table("picturedata").Where("pid = ?", pid).Update("isdelete", "1").Error; err != nil {
		return 0, ""
	}
	if err := db.Select("plast").Table("picturedata").Where("pid = ?", pid).First(&photoData).Error; err != nil {
		return 0, ""
	}
	return 1, photoData.Plast
}

/**
 * @Descripttion: 后台登录日志记录 - Back-end login log
 * @Author: William Wu
 * @Date: 2022/07/16 上午 09:51
 * @Param: uid (int)
 * @Param: account (string)
 * @Param: username (string)
 * @Param: userIp (string)
 */
func LogUserLogin(uid int, account string, username string, userIp string) {
	logData := AdminModel.UserLog{
		Uid:         uid,
		Account:     account,
		Username:    username,
		Ip:          userIp,
		Isdelete:    0,
		Create_Time: time.Now(),
		Update_Time: time.Now(),
	}
	isDeleteId := isDeleteLoginLogs()
	if isDeleteId != 0 {
		logData.Id = isDeleteId
		db.Table("login_logs").Save(&logData)
	} else {
		db.Table("login_logs").Create(&logData)
	}
}

/**
 * @Descripttion: 查询已被删除的日志 - Query isdelete log
 * @Author: William Wu
 * @Date: 2022/07/16 上午 09:55
 * @Return: Id (int)
 */
func isDeleteLoginLogs() int {
	logData := AdminModel.UserLog{}
	db.Select("id").Table("login_logs").First(&logData, "isdelete = ?", 1)
	return logData.Id
}

/**
 * @Descripttion: 清空后台登录日志 - clean back-end login log
 * @Author: William Wu
 * @Date: 2022/07/16 上午 10:41
 * @Return: Result (int)
 */
func CleanLoginLogs(uid int) int {
	if err := db.Table("login_logs").Where("uid = ? AND isdelete = ?", uid, 0).Update("isdelete", 1).Error; err != nil {
		return 0
	}
	return 1
}

/**
 * @Descripttion: 获取管理员的登录日志 - Get admin user login logs
 * @Author: William Wu
 * @Date: 2022/07/16 下午 12:35
 * @Param:
 * @Return:
 */
func GetLoginLogs(uid int) []AdminModel.UserLog {
	var logData []AdminModel.UserLog
	db.Table("login_logs").Where("uid = ? AND isdelete = ?", uid, 0).Limit(10).Find(&logData)
	return logData
}
