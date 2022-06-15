/*
 * @Descripttion: 视频观看控制层 - Watch Controller
 * @Author: William Wu
 * @Date: 2022/06/10 下午 02:37
 * @LastEditors: William Wu
 * @LastEditTime: 2022/06/10 下午 02:37
 */
package WatchController

import (
	"VideoHubGo/middlewares/JwtMiddleware"
	"VideoHubGo/models/WatchModel"
	"VideoHubGo/services/WatchServices"
	"VideoHubGo/utils/JsonUtils"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
 * @Descripttion: 查询视频详细信息 - Query video details
 * @Author: William Wu
 * @Date: 2022/06/10 下午 02:55
 * @Param: vid (int)
 * @Return: Json
 */
func GetVideoDetail(ctx *gin.Context) {
	requestBody := WatchModel.WatchRequest{}
	err := ctx.BindJSON(&requestBody)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "600", "参数错误 - Parameter error"))
		return
	}
	vid := requestBody.Vid
	uid := JwtMiddleware.GetTokenUID(ctx)
	videoData := WatchServices.GetVideoDetail(vid)
	reId := WatchServices.IsRelation(uid, vid)
	rData := map[string]interface{}{"list": videoData, "relation": reId}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", rData))
}

/**
 * @Descripttion: 增加视频流量 - Increase video traffic
 * @Author: William Wu
 * @Date: 2022/06/10 下午 02:55
 * @Param: vid (int)
 * @Return: Json
 */
func PlusVideoWatch(ctx *gin.Context) {
	requestBody := WatchModel.WatchRequest{}
	err := ctx.BindJSON(&requestBody)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "600", "参数错误 - Parameter error"))
		return
	}
	vid := requestBody.Vid
	result := WatchServices.PlusVideoWatch(vid)
	if result != 1 {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "增加视频流量时产生了异常 - Exception occurred while increasing video traffic\n\n", ""))
	}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", ""))
}
