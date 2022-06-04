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
 * @Return: 数据条数 - size (int)
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
	offset := size * (page - 1)

	var videoData = VideoCache.VideoGetListCache(page, size)

	if videoData == nil {
		//查询数据库数据 - Find Sql Data
		videoData = VideoServices.FindVideoList(size, offset)
		//缓存到Redis里 - Cache Redis
		VideoCache.VideoWriteListCache(videoData)
	}

	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", videoData))
}
