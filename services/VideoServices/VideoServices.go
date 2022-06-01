/*
 * @Descripttion: 视频服务层 - Video Services
 * @Author: William Wu
 * @Date: 2022/05/29 下午 08:38
 * @LastEditors: William Wu
 * @LastEditTime: 2022/05/29 下午 08:38
 */
package VideoServices

import (
	"VideoHubGo/models/VideoModel"
	"VideoHubGo/utils/DataBaseUtils"
)

var db = DataBaseUtils.GoDB()

/**
 * @Descripttion: 数据库查找视频数据 - Sql Find Video Data
 * @Author: William Wu
 * @Date: 2022/05/29 下午 08:46
 * @Param: size (int)
 * @Param: offset (int)
 * @Return: VideoModel Video
 */
func FindVideoList(size int, offset int) []VideoModel.VideoRe {
	var videoData []VideoModel.VideoRe
	db.Table("videodata").Where("isdelete = ?", 0).Order("create_time ASC").Limit(size).Offset(offset).Find(&videoData)
	return videoData
}
