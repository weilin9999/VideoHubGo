/*
 * @Descripttion: 用户收藏服务层 - Relation Services
 * @Author: William Wu
 * @Date: 2022/06/09 下午 03:27
 * @LastEditors: William Wu
 * @LastEditTime: 2022/06/09 下午 03:27
 */
package RelationServices

import (
	"VideoHubGo/models/VideoModel"
	"VideoHubGo/utils/DataBaseUtils"
)

var db = DataBaseUtils.GoDB()

/**
 * @Descripttion: 查询用户收藏信息集合 - Query user favorite information collection
 * @Author: William Wu
 * @Date: 2022/06/09 下午 03:29
 * @Param: uid (int)
 * @Param: size (int)
 * @Param: offset (int)
 * @Return: VideoModel Video
 */
func FindRelationByVideoList(uid int, size int, offset int) ([]VideoModel.VideoRe, int) {
	var videoData []VideoModel.VideoRe
	var count int64
	db.Select("videodata.vid,videodata.detail,videodata.watch,videodata.vtime,videodata.cid,videodata.create_time").Table("relation").Joins("LEFT JOIN videodata ON relation.vid = videodata.vid").Where(" relation.uid = ? AND relation.isdelete = ? AND videodata.isdelete = ?", uid, 0, 0).Order("videodata.vid DESC").Limit(size).Offset(offset).Count(&count).Find(&videoData)
	return videoData, int(count)
}
