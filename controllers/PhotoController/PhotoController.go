/*
 * @Descripttion: 图片控制层 Photo Controller
 * @Author: William Wu
 * @Date: 2022/07/08 下午 10:40
 * @LastEditors: William Wu
 * @LastEditTime: 2022/07/08 下午 10:40
 */
package PhotoController

import (
	"VideoHubGo/caches/PhotoCache"
	"VideoHubGo/models/PhotoModel"
	"VideoHubGo/services/PhotoServices"
	"VideoHubGo/utils/JsonUtils"
	"VideoHubGo/utils/LogUtils"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
	"net/http"
)

//使用SingleFlight防止Redis缓存击穿
var sfg singleflight.Group

/**
 * @Descripttion: 获取图片数据集 - Get photo data list
 * @Author: William Wu
 * @Date: 2022/07/08 下午 11:27
 * @Param: PhotoModel PhotoSearch (struct)
 * @Return: Json
 */
func GetPhotoList(ctx *gin.Context) {
	photoSearchData := PhotoModel.PhotoSearch{}
	err := ctx.BindJSON(&photoSearchData)
	if err != nil {
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(600, "参数错误 - Parameter error", ""))
		return
	}
	page := photoSearchData.Page
	size := photoSearchData.Size
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

	retData, err, _ := sfg.Do("getphotolist", func() (interface{}, error) {
		photoData := PhotoCache.PhotoGetListCache(page, size)
		count := PhotoCache.PhotoGetCount()
		redcount := PhotoCache.GetReidsPhotoListCount()

		if photoData == nil && len(photoData) < 20 {
			count = PhotoServices.GetCountPhotoList()
			PhotoCache.PhotoSaveCountList(count)
			if redcount < count {
				//查询数据库数据 - Find Sql Data
				photoData = PhotoServices.GetPhotoList(size, offset)
				//缓存到Redis里 - Cache Redis
				PhotoCache.PhotoWriteListCache(photoData)
			}
		}

		rData := map[string]interface{}{"list": photoData, "count": count}
		return rData, nil
	})
	if err != nil {
		LogUtils.Logger("[GIN ERROR]在处理getphotolist缓存击穿中异常：" + err.Error())
		ctx.JSON(http.StatusOK, JsonUtils.JsonResult(201, "未知的错误-击穿", ""))
		return
	}

	ctx.JSON(http.StatusOK, JsonUtils.JsonResult(200, "200", retData))
}
