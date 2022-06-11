/*
 * @Descripttion: 视频控制层 - Video Controller
 * @Author: William Wu
 * @Date: 2022/05/29 下午 05:05
 * @LastEditors: William Wu
 * @LastEditTime: 2022/05/29 下午 05:05
 */
package VideoController

import (
	"VideoHubGo/cache/VideoCache"
	"VideoHubGo/models/VideoModel"
	"VideoHubGo/services/VideoServices"
	"VideoHubGo/utils/JsonUtils"
	"github.com/gin-gonic/gin"
	"net/http"
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
	ctx.BindJSON(&requestBody)

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
	redcount := VideoCache.GetReidsVideoListCount()

	if videoData == nil && count != redcount {
		//查询数据库数据 - Find Sql Data
		videoData = VideoServices.FindVideoList(size, offset)
		//缓存到Redis里 - Cache Redis
		VideoCache.VideoWriteListCache(videoData)
		count := VideoServices.GetCountVideoList()
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
	ctx.BindJSON(&requestBody)
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
	ctx.BindJSON(&requestBody)
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
