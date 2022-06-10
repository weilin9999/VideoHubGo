/*
 * @Descripttion: 用户收藏控制层 - Relation Controller
 * @Author: William Wu
 * @Date: 2022/06/09 下午 03:20
 * @LastEditors: William Wu
 * @LastEditTime: 2022/06/09 下午 03:20
 */
package RelationController

import (
	"VideoHubGo/middlewares/JwtMiddleware"
	"VideoHubGo/models/RelationModel"
	"VideoHubGo/services/RelationServices"
	"VideoHubGo/utils/JsonUtils"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
 * @Descripttion: 获取用户收藏信息 - Get user favorite information
 * @Author: William Wu
 * @Date: 2022/06/09 下午 03:22
 * @Param: 分页 - page (int)
 * @Param: 数据条数 - size (int)
 * @Return: Json
 */
func GetRelationList(ctx *gin.Context) {
	requestBody := RelationModel.RelationRequest{}
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
	uid := JwtMiddleware.GetTokenUID(ctx)
	videoData, count := RelationServices.FindRelationByVideoList(uid, size, offset)
	rData := map[string]interface{}{"list": videoData, "count": count}
	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", rData))
}
