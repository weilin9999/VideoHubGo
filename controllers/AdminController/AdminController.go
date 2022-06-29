/*
 * @Descripttion: 管理员控制层 - Admin Controller
 * @Author: William Wu
 * @Date: 2022-05-21 11:51:55
 * @LastEditors: William Wu
 * @LastEditTime: 2022-05-21 18:35:03
 */
package AdminController

import (
	"VideoHubGo/utils/FileUtils"
	"VideoHubGo/utils/SystemUtils"
	"VideoHubGo/utils/UploadUtils"
	"github.com/gin-gonic/gin"
)

func AdminLogin(ctx *gin.Context) {
	//username := JwtMiddleware.GetTokenUsername(ctx)
	//dta := map[string]interface{}{
	//	"tag":      "<br>",
	//	"username": username,
	//}
	//ctx.AsciiJSON(200, dta)
	//filePath := UploadUtils.GetFilePath("video.saveFile") + "1571.mp4"
	filePath := UploadUtils.GetFilePath("video.saveFile")
	//coverPath := UploadUtils.GetFilePath("video.coverFile") + "1571"
	//row := FFmpegUtils.GetVideoCover(filePath, coverPath)
	//if row != 1 {
	//	ctx.AsciiJSON(200, "封面截取失败")
	//	return
	//}
	//row2 := FFmpegUtils.RecodeVideo(filePath, filePath+"1571.mp4")
	//if row2 != 1 {
	//	ctx.AsciiJSON(200, "封面截取失败")
	//	return
	//}
	//res, _ := VideoUtils.GetVideoTotalTime(filePath)
	res, _ := FileUtils.TraverseFile(filePath)
	disk, _ := SystemUtils.GetDiskInfo(filePath)
	dta := map[string]interface{}{
		"tag":  res,
		"disk": disk,
	}
	ctx.AsciiJSON(200, dta)
}
