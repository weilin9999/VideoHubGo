/*
 * @Descripttion: 视频控制层 - Video Controller
 * @Author: William Wu
 * @Date: 2022/05/29 下午 05:05
 * @LastEditors: William Wu
 * @LastEditTime: 2022/05/29 下午 05:05
 */
package VideoController

import (
	"VideoHubGo/caches/VideoCache"
	"VideoHubGo/middlewares/JwtMiddleware"
	"VideoHubGo/models/VideoModel"
	"VideoHubGo/services/VideoServices"
	"VideoHubGo/utils/JsonUtils"
	"VideoHubGo/utils/UploadUtils"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
)

/**
 * @Descripttion: 获取视频数据 - Get Video List Data
 * @Author: William Wu
 * @Date: 2022/05/29 下午 08:54
 * @Param: 分页 - page (int)
 * @Param: 数据条数 - size (int)
 * @Return: Json
 */
func GetVideoList(ctx *gin.Context) {
	requestBody := VideoModel.VideoRequest{}
	err := ctx.BindJSON(&requestBody)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}

	page := requestBody.Page
	size := requestBody.Size
	if page < 1 {
		page = 1
	}
	if size < 20 {
		size = 20
	}
	if size > 40 {
		size = 20
	}

	offset := size * (page - 1)

	var videoData = VideoCache.VideoGetListCache(page, size)
	count := VideoCache.VideoGetCount()
	//redcount := VideoCache.GetReidsVideoListCount()

	if videoData == nil {
		//查询数据库数据 - Find Sql Data
		videoData = VideoServices.FindVideoList(size, offset)
		//缓存到Redis里 - Cache Redis
		VideoCache.VideoWriteListCache(videoData)
		count = VideoServices.GetCountVideoList()
		VideoCache.VideoSaveCountList(count)
	}

	rData := map[string]interface{}{"list": videoData, "count": count}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", rData))
}

/**
 * @Descripttion: 获取定义分类的视频数据 - Get video data defining classification
 * @Author: William Wu
 * @Date: 2022/06/05 下午 06:16
 * @Param: cid (int)
 * @Param: 分页 - page (int)
 * @Param: 数据条数 - size (int)
 * @Return: Json
 */
func GetVideoClassList(ctx *gin.Context) {
	requestBody := VideoModel.VideoRequestClass{}
	err := ctx.BindJSON(&requestBody)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	cid := requestBody.Cid
	page := requestBody.Page
	size := requestBody.Size
	if page < 1 {
		page = 1
	}
	if size < 20 {
		size = 20
	}
	if size > 40 {
		size = 20
	}
	if cid < 0 {
		cid = 0
	}
	offset := size * (page - 1)
	count := VideoCache.VideoGetClassCount(cid)
	if count == 0 {
		tempCount := VideoServices.GetCountVideoClassList(cid)
		VideoCache.VideoSaveClassCountList(cid, tempCount)
		count = tempCount
	}
	var videoData = VideoServices.FindVideoInClass(cid, size, offset)
	rData := map[string]interface{}{"list": videoData, "count": count}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", rData))
}

/**
 * @Descripttion: 搜索视频 - Search Video
 * @Author: William Wu
 * @Date: 2022/06/08 下午 02:29
 * @Param: 分类ID - cid (int)
 * @Param: 搜搜关键字 - key (string)
 * @Param: 分页 - page (int)
 * @Param: 数据条数 - size (int)
 * @Return: Json
 */
func SearchVideoClassList(ctx *gin.Context) {
	requestBody := VideoModel.VideoRequestSearch{}
	err := ctx.BindJSON(&requestBody)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	cid := requestBody.Cid
	key := requestBody.Key
	page := requestBody.Page
	size := requestBody.Size
	if page < 1 {
		page = 1
	}
	if size < 20 {
		size = 20
	}
	if size > 40 {
		size = 20
	}
	if cid < 0 {
		cid = 0
	}
	offset := size * (page - 1)
	var videoData, count = VideoServices.SearchVideoList_Class(cid, key, size, offset)
	rData := map[string]interface{}{"list": videoData, "count": count}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", rData))
}

/**
 * @Descripttion: 视频上传 - Video Upload
 * @Author: William Wu
 * @Date: 2022/06/15 上午 11:09
 * @Param: File
 * @Return: Json
 */
func UploadVideo_StreamFile(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	detail := ctx.PostForm("detail")
	cid := ctx.PostForm("cid")
	recid, err := strconv.Atoi(cid)
	uid := JwtMiddleware.GetTokenUID(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(602, "cid基本参数错误 - CID Base Parameter error", ""))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "文件参数错误 - File Parameter error", ""))
		return
	}
	if len(detail) < 1 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(601, "detail基本参数错误 - Detail Base Parameter error", ""))
		return
	}
	savePath := UploadUtils.GetFilePath("video.saveFile")
	fileSuffix := path.Ext(header.Filename)
	if fileSuffix != ".mp4" {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "文件类型不符合要求 - Document type does not meet requirements", ""))
		return
	}
	saveVid := VideoServices.UploadVideo(uid, detail, recid)
	save, err := os.OpenFile(savePath+strconv.Itoa(saveVid)+fileSuffix, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(202, "读取文件流失败 - Failed to read file stream", ""))
		return
	}
	for {
		buf := make([]byte, 1024*2)
		read, err := file.Read(buf)
		if err != nil && err != io.EOF {
			ctx.JSON(http.StatusOK, JsonUtils.JsonResult(203, "视频上传出现异常 - Abnormal video uploading", ""))
			return
		}
		if read == 0 {
			break
		}
		_, err = save.Write(buf)
		if err != nil {
			ctx.JSON(http.StatusOK, JsonUtils.JsonResult(204, "视频存储过程出现异常 - Exception in video stored procedure", ""))
			return
		}
	}
	rData := map[string]interface{}{"vid": saveVid}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "视频上传成功 - Video upload succeeded", rData))
}
