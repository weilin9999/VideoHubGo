/*
 * @Descripttion: 后台控制层 - Admin Controller
 * @Author: William Wu
 * @Date: 2022-05-21 11:51:55
 * @LastEditors: William Wu
 * @LastEditTime: 2022-05-21 18:35:03
 */
package AdminController

import (
	"VideoHubGo/caches/AdminCache"
	"VideoHubGo/caches/ClassCache"
	"VideoHubGo/caches/DistributedLock"
	"VideoHubGo/caches/PhotoCache"
	"VideoHubGo/caches/VideoCache"
	"VideoHubGo/middlewares/JwtMiddleware"
	"VideoHubGo/models/AdminModel"
	"VideoHubGo/models/ClassModel"
	"VideoHubGo/models/UserModel"
	"VideoHubGo/models/VideoModel"
	"VideoHubGo/services/AdminServices"
	"VideoHubGo/services/PhotoServices"
	"VideoHubGo/services/UserServices"
	"VideoHubGo/services/VideoServices"
	"VideoHubGo/utils/FFmpegUtils"
	"VideoHubGo/utils/FileUtils"
	"VideoHubGo/utils/JsonUtils"
	"VideoHubGo/utils/LogUtils"
	"VideoHubGo/utils/SystemUtils"
	"VideoHubGo/utils/UploadUtils"
	"VideoHubGo/utils/VideoUtils"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
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

	userData = AdminServices.LoginAdmin(userData.Account, userData.Password)

	if userData.Uid == 0 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(202, "无效会话 - TimeOut", ""))
		return
	}
	token, err := JwtMiddleware.CreateToken(userData.Uid, userData.Account, userData.Username, userData.Isadmin, userData.Isuploader)

	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(205, "Token生成失败 - Token Create Error", ""))
		return
	}

	userIp := ctx.ClientIP()
	AdminServices.LogUserLogin(userData.Uid, userData.Account, userData.Username, userIp)

	rData := map[string]interface{}{"token": token, "username": userData.Username, "avatar": userData.Avatar, "account": userData.Account}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", rData))
}

/**
 * @Descripttion: 后台搜索视频 - Admin search video
 * @Author: William Wu
 * @Date: 2022/07/05 下午 01:03
 * @Param: VideoModel.VideoAdminSearch (struct)
 * @Return: Json
 */
func GetVideoList(ctx *gin.Context) {
	videoSearchData := VideoModel.VideoAdminSearch{}
	err := ctx.BindJSON(&videoSearchData)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	vid := videoSearchData.Vid
	cid := videoSearchData.Cid
	detail := videoSearchData.Detail
	page := videoSearchData.Page
	size := videoSearchData.Size
	if page < 1 {
		page = 1
	}
	if size < 10 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	offset := size * (page - 1)
	videoData, count := AdminServices.GetVideoList(vid, detail, cid, size, offset)
	rData := map[string]interface{}{"list": videoData, "count": count}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", rData))
}

/**
 * @Descripttion: 搜索未分类的视频 - Seach no cid video data
 * @Author: William Wu
 * @Date: 2022/07/08 下午 06:05
 * @Param: VideoModel.VideoAdminSearch (struct)
 * @Return: Json
 */
func GetNoCidVideoList(ctx *gin.Context) {
	videoSearchData := VideoModel.VideoAdminSearch{}
	err := ctx.BindJSON(&videoSearchData)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	page := videoSearchData.Page
	size := videoSearchData.Size
	if page < 1 {
		page = 1
	}
	if size < 10 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	offset := size * (page - 1)
	videoData, count := AdminServices.GetNoCidVideoList(size, offset)
	rData := map[string]interface{}{"list": videoData, "count": count}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", rData))
}

/**
 * @Descripttion: 批量删除视频信息 - Batch delete video information
 * @Author: William Wu
 * @Date: 2022/07/05 下午 07:45
 * @Param: vid (Array)
 * @Return: Json
 */
func DeleteVideoGroup(ctx *gin.Context) {
	delteMap := AdminModel.VideoDeleteGroup{}
	err := ctx.BindJSON(&delteMap)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	newFile := UploadUtils.GetFilePath("waste.file")
	oldFile := UploadUtils.GetFilePath("video.saveFile")
	coverPath := UploadUtils.GetFilePath("video.coverFile")
	errCount := 0
	for _, videoDelete := range delteMap.Group {
		res, name := AdminServices.DeleteVideoInformation(videoDelete.Vid)
		if res != 1 {
			errCount++
		} else {
			FileUtils.MoveFile(oldFile+strconv.Itoa(videoDelete.Vid)+".mp4", newFile+name+".mp4")
			FileUtils.DeleteFile(coverPath + strconv.Itoa(videoDelete.Vid) + ".png")
		}
	}
	if errCount != 0 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "批量删除途中出现异常", ""))
		return
	}
	VideoCache.VideoDeleteCaches()
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", ""))
}

/**
 * @Descripttion: 修改视频信息 - Update video information
 * @Author: William Wu
 * @Date: 2022/06/30 下午 08:27
 * @Param: VideoEdit (struct)
 * @Return: Josn
 */
func EditVideoInformation(ctx *gin.Context) {
	videoData := VideoModel.VideoEdit{}
	err := ctx.BindJSON(&videoData)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	if res := AdminServices.EditVideoInformation(videoData); res == 0 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "修改失败 - Update fail", ""))
	} else {
		VideoCache.VideoDeleteCaches()
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "修改成功 - Update success", ""))
	}
}

/**
 * @Descripttion: 批量修改视频的分类 - Batch edit video cid information
 * @Author: William Wu
 * @Date: 2022/07/05 下午 08:32
 * @Param: VideoBatchCid (struct)
 * @Return: Josn
 */
func BatchVideoCidInformation(ctx *gin.Context) {
	batchMap := AdminModel.VideoBatchCid{}
	err := ctx.BindJSON(&batchMap)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	errCount := 0
	for _, videoDelete := range batchMap.Group {
		res := AdminServices.EditVideoCidInformation(batchMap.Cid, videoDelete.Vid)
		if res != 1 {
			errCount++
		}
	}
	if errCount != 0 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "批量修改视频分类途中出现异常", ""))
		return
	}
	VideoCache.VideoDeleteCaches()
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", ""))
}

/**
 * @Descripttion: 删除视频 - Delete Video
 * @Author: William Wu
 * @Date: 2022/06/30 下午 09:52
 * @Param: VideoEdit (struct) - vid (int)
 * @Return: Josn
 */
func DeleteVideoInformation(ctx *gin.Context) {
	videoData := VideoModel.VideoEdit{}
	err := ctx.BindJSON(&videoData)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	res, name := AdminServices.DeleteVideoInformation(videoData.Vid)
	if res == 0 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "删除失败 - Delete fail", ""))
	} else {
		newFile := UploadUtils.GetFilePath("waste.file")
		oldFile := UploadUtils.GetFilePath("video.saveFile")
		coverPath := UploadUtils.GetFilePath("video.coverFile")
		FileUtils.MoveFile(oldFile+strconv.Itoa(videoData.Vid)+".mp4", newFile+name+".mp4")
		FileUtils.DeleteFile(coverPath + strconv.Itoa(videoData.Vid) + ".png")
		VideoCache.VideoDeleteCaches()
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "删除成功 - Delete success", ""))
	}
}

/**
 * @Descripttion: 后台获取用户数据 - Admin get user data
 * @Author: William Wu
 * @Date: 2022/07/05 下午 11:02
 * @Param: UserAdminSearch (struct)
 * @Return: Json
 */
func GetUserList(ctx *gin.Context) {
	userSearchData := UserModel.UserAdminSearch{}
	err := ctx.BindJSON(&userSearchData)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	uid := userSearchData.Uid
	account := userSearchData.Account
	username := userSearchData.Username
	page := userSearchData.Page
	size := userSearchData.Size
	if page < 1 {
		page = 1
	}
	if size < 10 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	offset := size * (page - 1)
	userData, count := AdminServices.GetUserList(uid, account, username, size, offset)
	rData := map[string]interface{}{"list": userData, "count": count}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", rData))
}

/**
 * @Descripttion: 修改用户信息 - Update User Information
 * @Author: William Wu
 * @Date: 2022/06/30 下午 10:22
 * @Param: EditUser (struct)
 * @Return: Josn
 */
func EditUserInformation(ctx *gin.Context) {
	userData := UserModel.EditUser{}
	err := ctx.BindJSON(&userData)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	if res := AdminServices.EditUserInformation(userData); res == 0 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "修改失败 - Update fail", ""))
	} else if res == 2 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(202, "修改密码失败 - Update password fail", ""))
	} else {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "修改成功 - Update success", ""))
	}
}

/**
 * @Descripttion: 提升/下降用户权限 - Promote or Decline User authority
 * @Author: William Wu
 * @Date: 2022/07/05 下午 11:17
 * @Param: uid (int)
 * @Param: isadmin (int)
 * @Return: Json
 */
func AuthorityUserIsadmin(ctx *gin.Context) {
	userData := AdminModel.UserIsadmin{}
	err := ctx.BindJSON(&userData)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	if res := AdminServices.AuthorityUserIsadmin(userData.Uid, userData.Isadmin); res == 0 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "修改失败 - Update fail", ""))
	} else {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "修改成功 - Update success", ""))
	}
}

/**
 * @Descripttion: 提升/下降用户权限 - Promote or Decline User authority
 * @Author: William Wu
 * @Date: 2022/07/05 下午 11:19
 * @Param: uid (int)
 * @Param: isuploader (int)
 * @Return: Json
 */
func AuthorityUserIsuploader(ctx *gin.Context) {
	userData := AdminModel.UserIsuploader{}
	err := ctx.BindJSON(&userData)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	if res := AdminServices.AuthorityUserIsuploader(userData.Uid, userData.Isuploader); res == 0 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "修改失败 - Update fail", ""))
	} else {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "修改成功 - Update success", ""))
	}
}

/**
 * @Descripttion:
 * @Author: William Wu
 * @Date: 2022/07/05 下午 11:08
 * @Param: ClassDeleteGroup (struct)
 * @Return: Josn
 */
func DeleteUserGroup(ctx *gin.Context) {
	delteMap := AdminModel.UserDeleteGroup{}
	err := ctx.BindJSON(&delteMap)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	errCount := 0
	for _, userDelete := range delteMap.Group {
		res := AdminServices.DeleteUserInformation(userDelete.Uid)
		if res != 1 {
			errCount++
		}
	}
	if errCount != 0 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "批量删除途中出现异常", ""))
		return
	}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", ""))
}

/**
 * @Descripttion: 删除用户信息 - Delete User Information
 * @Author: William Wu
 * @Date: 2022/06/30 下午 11:12
 * @Param: EditUser (struct) - uid (int)
 * @Return: Josn
 */
func DeleteUserInformation(ctx *gin.Context) {
	userData := UserModel.EditUser{}
	err := ctx.BindJSON(&userData)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	if res := AdminServices.DeleteUserInformation(userData.Uid); res == 0 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "删除失败 - Delete fail", ""))
	} else {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "删除成功 - Delete success", ""))
	}
}

/**
 * @Descripttion: 后台获取分类信息 - Admin get class informarion
 * @Author: William Wu
 * @Date: 2022/07/05 下午 09:23
 * @Param: ClassModel.ClassAdminSearch (struct)
 * @Return: Json
 */
func GetClassList(ctx *gin.Context) {
	classSearchData := ClassModel.ClassAdminSearch{}
	err := ctx.BindJSON(&classSearchData)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	cid := classSearchData.Cid
	classname := classSearchData.Classname
	page := classSearchData.Page
	size := classSearchData.Size
	if page < 1 {
		page = 1
	}
	if size < 10 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	offset := size * (page - 1)
	classData, count := AdminServices.GetClassList(cid, classname, size, offset)
	rData := map[string]interface{}{"list": classData, "count": count}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", rData))
}

/**
 * @Descripttion: 创建新视频分类 - Created new video class
 * @Author: William Wu
 * @Date: 2022/07/09 上午 11:12
 * @Param: classname (string)
 * @Return: Json
 */
func CreatedClassInformarion(ctx *gin.Context) {
	classData := ClassModel.ClassRe{}
	err := ctx.BindJSON(&classData)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	res := AdminServices.CreatedClassInformarion(classData.Classname)
	if res == 2 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "在操作插入新视频分类更新时产生异常", ""))
	} else if res == 3 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(202, "在操作插入新视频分类创建时产生异常", ""))
	} else {
		ClassCache.ClassDeleteCaches()
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "创建新视频分类成功", ""))
	}
}

/**
 * @Descripttion: 修改视频分类信息 - Update Video Class Information
 * @Author: William Wu
 * @Date: 2022/06/30 下午 11:27
 * @Param: ClassRe (struct)
 * @Return: Josn
 */
func EditClassInformation(ctx *gin.Context) {
	classData := ClassModel.ClassRe{}
	err := ctx.BindJSON(&classData)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	if res := AdminServices.EditClassInformation(classData); res == 0 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "修改失败 - Update fail", ""))
	} else {
		ClassCache.ClassDeleteCaches()
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "修改成功 - Update success", ""))
	}
}

/**
 * @Descripttion: 批量删除分类信息 - Batch delete class information
 * @Author: William Wu
 * @Date: 2022/07/05 下午 09:27
 * @Param: ClassDeleteGroup (struct)
 * @Return: Josn
 */
func DeleteClassGroup(ctx *gin.Context) {
	delteMap := AdminModel.ClassDeleteGroup{}
	err := ctx.BindJSON(&delteMap)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	errCount := 0
	for _, classDelete := range delteMap.Group {
		res := AdminServices.DeleteClassInformation(classDelete.Cid)
		if res != 1 {
			errCount++
		}
	}
	if errCount != 0 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "批量删除途中出现异常", ""))
		return
	}
	ClassCache.ClassDeleteCaches()
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", ""))
}

/**
 * @Descripttion: 删除视频分类信息 - Delete Video Class Information
 * @Author: William Wu
 * @Date: 2022/06/30 下午 11:28
 * @Param: ClassRe (struct) - cid (int)
 * @Return: Josn
 */
func DeleteClassInformation(ctx *gin.Context) {
	classData := ClassModel.ClassRe{}
	err := ctx.BindJSON(&classData)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	if res := AdminServices.DeleteClassInformation(classData.Cid); res == 0 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "删除失败 - Delete fail", ""))
	} else {
		ClassCache.ClassDeleteCaches()
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "删除成功 - Delete success", ""))
	}
}

/**
 * @Descripttion: 后台更改用户头像 - Admin edit user avatar
 * @Author: William Wu
 * @Date: 2022/07/06 上午 11:09
 * @Param: uid (int)
 * @Param: file (File)
 * @Return: Json
 */
func UploadAvatar(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	uid := ctx.PostForm("uid")
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
			ctx.JSON(http.StatusOK, JsonUtils.JsonResult(203, "文件类型不符合要求 - Document type does not meet requirements", ""))
			return
		}
		userID, _ := strconv.Atoi(uid)
		fileName := UserServices.UserUploadAvatar(userID, fileType)
		filePath := filepath.Join(UploadUtils.GetFilePath("user.userAvatar"), fileName)
		err := ctx.SaveUploadedFile(file, filePath)
		if err != nil {
			ctx.JSON(http.StatusOK, JsonUtils.JsonResult(204, "在存储头像时出现了异常 - An exception occurred while storing the Avatar", ""))
			LogUtils.Logger("异常错误-Error：头像保存处理时产生 - Generated when saving the Avatar ：" + err.Error())
			return
		}
		rData := map[string]interface{}{"avatar": fileName}
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "上传头像成功 - Upload Avatar Success", rData))
	}
}

/**
 * @Descripttion: 获取后台服务器的存储空间信息 - Get back-end server strong information
 * @Author: William Wu
 * @Date: 2022/07/06 下午 04:27
 * @Return: Json
 */
func GetDiskStrong(ctx *gin.Context) {
	strong, _ := SystemUtils.GetDiskInfo(UploadUtils.GetFilePath("video.saveFile"))
	rData := map[string]interface{}{"strong": strong}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", rData))
}

/**
 * @Descripttion: 获取二次编码状态 - Get recode status
 * @Author: William Wu
 * @Date: 2022/07/06 下午 06:32
 * @Return: ReCode Status (int)
 */
func GetReCodeStatus(ctx *gin.Context) {
	cacheStatus := AdminCache.GetReCodeCache()
	if cacheStatus == 0 {
		cacheStatus = AdminServices.GetReCodeStatus()
		AdminCache.SaveReCodeCountList(cacheStatus)
	}
	rData := map[string]interface{}{"recode": cacheStatus}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", rData))
}

/**
 * @Descripttion: 修改后台二次编码状态 - Change back-end recode status
 * @Author: William Wu
 * @Date: 2022/07/06 下午 07:03
 * @Param: status (int)
 * @Return: Json
 */
func ChangeReCodeStatus(ctx *gin.Context) {
	reCodeData := AdminModel.AdminDashboard{}
	err := ctx.BindJSON(&reCodeData)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	if res := AdminServices.ChangeReCodeStatus(reCodeData.Re_Code); res == 0 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "修改失败 - Change fail", ""))
	} else {
		AdminCache.ReCodeDeleteCaches()
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "修改成功 - Change success", ""))
	}
}

/**
 * @Descripttion: 获取后台临时文件列表 - Get admin temp file list
 * @Author: William Wu
 * @Date: 2022/07/07 下午 01:28
 */
func GetTempFileList(ctx *gin.Context) {
	fileList, err := FileUtils.TraverseFile(UploadUtils.GetFilePath("video.tempFile"))
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "error", ""))
		return
	}
	rData := map[string]interface{}{"list": fileList}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", rData))
}

/**
 * @Descripttion: 获取后台回收站文件列表 - Get admin recycle bin file list
 * @Author: William Wu
 * @Date: 2022/07/07 下午 09:56
 */
func GetWasteFileList(ctx *gin.Context) {
	fileList, err := FileUtils.TraverseFile(UploadUtils.GetFilePath("waste.file"))
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "error", ""))
		return
	}
	rData := map[string]interface{}{"list": fileList}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", rData))
}

/**
 * @Descripttion: 删除临时文件的文件到回收站 - Delete temp file to recycle bin
 * @Author: William Wu
 * @Date: 2022/07/07 下午 10:28
 * @Param: name (string)
 * @Return: Json
 */
func DeleteTempFile(ctx *gin.Context) {
	deleteData := AdminModel.DeleteStruct{}
	err := ctx.BindJSON(&deleteData)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	oldFile := UploadUtils.GetFilePath("video.tempFile")
	newFile := UploadUtils.GetFilePath("waste.file")
	res := FileUtils.MoveFile(oldFile+deleteData.Name, newFile+deleteData.Name)
	if res == 0 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "删除文件失败", ""))
		return
	}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "删除文件成功", ""))
}

/**
 * @Descripttion: 删除回收站文件 - Delete recycle bin file
 * @Author: William Wu
 * @Date: 2022/07/07 下午 10:31
 * @Param: name (string)
 * @Return: Json
 */
func DeleteWasteFile(ctx *gin.Context) {
	deleteData := AdminModel.DeleteStruct{}
	err := ctx.BindJSON(&deleteData)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	deleFile := UploadUtils.GetFilePath("waste.file")
	res := FileUtils.DeleteFile(deleFile + deleteData.Name)
	if res == 0 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "删除文件失败", ""))
		return
	}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "删除文件成功", ""))
}

/**
 * @Descripttion: 恢复回收站文件 - Recovery recycle bin file
 * @Author: William Wu
 * @Date: 2022/07/07 下午 11:45
 * @Param: name (string)
 * @Return: Json
 */
func RecoveryWasteFile(ctx *gin.Context) {
	recoveryData := AdminModel.DeleteStruct{}
	err := ctx.BindJSON(&recoveryData)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	oldFile := UploadUtils.GetFilePath("waste.file")
	newFile := UploadUtils.GetFilePath("video.tempFile")
	res := FileUtils.MoveFile(oldFile+recoveryData.Name, newFile+recoveryData.Name)
	if res == 0 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "恢复文件失败", ""))
		return
	}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "恢复文件成功", ""))
}

/**
 * @Descripttion: 后台临时文件上传 - Back-end upload temp file
 * @Author: William Wu
 * @Date: 2022/07/07 下午 07:50
 * @Param: file (file)
 * @Return: Json
 */
func UploadFileStream(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "文件参数错误 - File Parameter error", ""))
		return
	}
	savePath := UploadUtils.GetFilePath("video.tempFile")
	save, err := os.OpenFile(savePath+header.Filename, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "读取文件流失败 - Failed to read file stream", ""))
		return
	}
	for {
		buf := make([]byte, 1024*2)
		read, err := file.Read(buf)
		if err != nil && err != io.EOF {
			ctx.JSON(http.StatusOK, JsonUtils.JsonResult(202, "文件上传出现异常 - Abnormal file uploading", ""))
			return
		}
		if read == 0 {
			break
		}
		_, err = save.Write(buf)
		if err != nil {
			ctx.JSON(http.StatusOK, JsonUtils.JsonResult(203, "文件存储过程出现异常 - Exception in file stored procedure", ""))
			return
		}
	}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "文件上传成功 - file upload succeeded", ""))
}

/**
 * @Descripttion: 将临时文件夹的视频存储到服务器中 - Temp file video save to server
 * @Author: William Wu
 * @Date: 2022/07/08 下午 01:57
 * @Return: Json
 */
func SaveVideoToLocal(ctx *gin.Context) {
	isLock := DistributedLock.GetLock("videolock")
	if isLock == 1 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(400, "后台已经在运行请勿重复操作", ""))
		return
	} else if isLock == 2 {
		DistributedLock.UnLockFunc("videolock")
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "视频存储完成", ""))
		return
	} else {
		DistributedLock.LockFunc("videolock", 1)
	}
	filePath := UploadUtils.GetFilePath("video.tempFile")
	uid := ctx.MustGet("uid").(int)
	ioRead, err := ioutil.ReadDir(filePath)
	if err != nil {
		LogUtils.Logger("[IO操作异常 - IO operation exception] 遍历文件时出现异常 Traverse file operation exception：" + err.Error())
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "服务器IO异常", ""))
		return
	}
	recodeStatus := AdminServices.GetReCodeStatus()
	savePath := UploadUtils.GetFilePath("video.saveFile")
	coverPath := UploadUtils.GetFilePath("video.coverFile")
	for _, file := range ioRead {
		fileSuffix := path.Ext(file.Name())
		//如果不是.mp4文件直接跳过此次循环
		if fileSuffix != ".mp4" {
			continue
		}
		fileName := file.Name()
		detail := fileName[0 : len(fileName)-len(fileSuffix)]
		videoTime, err := VideoUtils.GetVideoTotalTime(filePath + fileName)
		if err != nil {
			continue
		}
		saveVid := VideoServices.UploadVideo(uid, detail, 0, videoTime)
		if recodeStatus == 2 {
			FFmpegUtils.RecodeVideo(filePath+fileName, savePath+strconv.Itoa(saveVid)+".mp4")
		} else {
			FileUtils.MoveFile(filePath+fileName, savePath+strconv.Itoa(saveVid)+".mp4")
		}
		FFmpegUtils.GetVideoCover(savePath+strconv.Itoa(saveVid)+".mp4", coverPath+strconv.Itoa(saveVid))
		FileUtils.DeleteFile(filePath + fileName)
	}
	VideoCache.VideoDeleteCaches()
	DistributedLock.LockFunc("videolock", 2)
}

/**
 * @Descripttion: 批量删除图片 - Batch delete photo information
 * @Author: William Wu
 * @Date: 2022/07/08 下午 11:51
 * @Param: PhotoDeleteGroup (struct)
 * @Return: Josn
 */
func DeletePhotoGroup(ctx *gin.Context) {
	delteMap := AdminModel.PhotoDeleteGroup{}
	err := ctx.BindJSON(&delteMap)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	errCount := 0
	for _, photoDelete := range delteMap.Group {
		res, plast := AdminServices.DeletePhotoInformation(photoDelete.Pid)
		if res != 1 {
			errCount++
		} else {
			newFile := UploadUtils.GetFilePath("video.tempFile")
			photoPath := UploadUtils.GetFilePath("photo.file")
			FileUtils.MoveFile(photoPath+strconv.Itoa(photoDelete.Pid)+plast, newFile+strconv.Itoa(photoDelete.Pid)+plast)
		}
	}
	if errCount != 0 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "批量删除途中出现异常", ""))
		return
	}
	PhotoCache.PhotoDeleteCaches()
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", ""))
}

/**
 * @Descripttion: 删除图片分类信息 - Delete Photo Class Information
 * @Author: William Wu
 * @Date: 2022/07/09 上午 12:05
 * @Param: PhotoDelete (struct) - pid (int)
 * @Return: Josn
 */
func DeletePhotoInformation(ctx *gin.Context) {
	photoData := AdminModel.PhotoDelete{}
	err := ctx.BindJSON(&photoData)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	res, plast := AdminServices.DeletePhotoInformation(photoData.Pid)
	if res == 0 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "删除失败 - Delete fail", ""))
	} else {
		newFile := UploadUtils.GetFilePath("video.tempFile")
		photoPath := UploadUtils.GetFilePath("photo.file")
		FileUtils.MoveFile(photoPath+strconv.Itoa(photoData.Pid)+plast, newFile+strconv.Itoa(photoData.Pid)+plast)
		PhotoCache.PhotoDeleteCaches()
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "删除成功 - Delete success", ""))
	}
}

/**
 * @Descripttion: 将临时文件夹的图片存储到服务器中 - Temp file photo save to server
 * @Author: William Wu
 * @Date: 2022/07/09 上午 01:19
 * @Return: Json
 */
func SavePhotoToLocal(ctx *gin.Context) {
	isLock := DistributedLock.GetLock("photolock")
	if isLock == 1 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(400, "后台已经在运行请勿重复操作", ""))
		return
	} else if isLock == 2 {
		DistributedLock.UnLockFunc("photolock")
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "图片存储完成", ""))
		return
	} else {
		DistributedLock.LockFunc("photolock", 1)
	}
	filePath := UploadUtils.GetFilePath("video.tempFile")
	ioRead, err := ioutil.ReadDir(filePath)
	if err != nil {
		LogUtils.Logger("[IO操作异常 - IO operation exception] Photo 遍历文件时出现异常 Traverse file operation exception：" + err.Error())
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "服务器IO异常", ""))
		return
	}
	savePath := UploadUtils.GetFilePath("photo.file")
	for _, file := range ioRead {
		fileSuffix := path.Ext(file.Name())
		//如果不是.mp4文件直接跳过此次循环
		if fileSuffix != ".jpg" && fileSuffix != ".jpeg" && fileSuffix != ".png" {
			continue
		}
		fileName := file.Name()
		saveVid := PhotoServices.SavePhoto(fileSuffix)
		FileUtils.MoveFile(filePath+fileName, savePath+strconv.Itoa(saveVid)+fileSuffix)
	}
	PhotoCache.PhotoDeleteCaches()
	DistributedLock.LockFunc("photolock", 2)
}

/**
 * @Descripttion: 清空后台登录日志 - clean back-end login log
 * @Author: William Wu
 * @Date: 2022/07/16 上午 10:43
 * @Return: Json
 */
func CleanLoginLogs(ctx *gin.Context) {
	uid := ctx.MustGet("uid").(int)
	res := AdminServices.CleanLoginLogs(uid)
	if res != 1 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "清空失败-请尝试重试", ""))
		return
	}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "已全部清除", ""))
}

/**
 * @Descripttion: 获取管理员的登录日志 - Get admin user login logs
 * @Author: William Wu
 * @Date: 2022/07/16 下午 12:39
 * @Return: Json
 */
func GetLoginLogs(ctx *gin.Context) {
	uid := ctx.MustGet("uid").(int)
	logData := AdminServices.GetLoginLogs(uid)
	rData := map[string]interface{}{"list": logData}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", rData))
}
