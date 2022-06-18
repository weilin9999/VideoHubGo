/*
 * @Descripttion: 视频分类控制器 - Video Class Controller
 * @Author: William Wu
 * @Date: 2022/06/03 下午 02:21
 * @LastEditors: William Wu
 * @LastEditTime: 2022/06/03 下午 02:21
 */
package ClassController

import (
	"VideoHubGo/caches/ClassCache"
	"VideoHubGo/services/ClassServices"
	"VideoHubGo/utils/JsonUtils"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
 * @Descripttion: 获取所有的视频分类 - Get All Video Class List
 * @Author: William Wu
 * @Date: 2022/06/03 下午 02:24
 * @Return: 数据 - Data
 */
func GetClassList(ctx *gin.Context) {
	var classData = ClassCache.ClassGetListCache()
	if classData == nil {
		classData = ClassServices.FindAllClass()
		//存入Redis - Save In Redis
		ClassCache.ClassWriteCache(classData)
	}

	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", classData))
}
