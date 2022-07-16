/*
 * @Descripttion: 视频观看服务层 - Watch Services
 * @Author: William Wu
 * @Date: 2022/06/10 下午 02:40
 * @LastEditors: William Wu
 * @LastEditTime: 2022/06/10 下午 02:40
 */
package WatchServices

import (
	"VideoHubGo/models/RelationModel"
	"VideoHubGo/models/VideoModel"
	"VideoHubGo/utils/DataBaseUtils"
	"VideoHubGo/utils/LogUtils"
	"gorm.io/gorm"
)

var db = DataBaseUtils.GoDB()

/**
 * @Descripttion: 查询视频详细内容 - Query video details
 * @Author: William Wu
 * @Date: 2022/06/10 下午 02:42
 * @Param: vid (int)
 * @Return: VideoModel VideoRe
 */
func GetVideoDetail(vid int) VideoModel.VideoRe {
	var videoData VideoModel.VideoRe
	db.Table("videodata").Where("isdelete = ? AND vid = ?", 0, vid).Take(&videoData)
	return videoData
}

/**
 * @Descripttion: 增加视频观看次数 - Increase video viewing
 * @Author: William Wu
 * @Date: 2022/06/10 下午 02:49
 * @Param: vid (int)
 * @Return: success (int)
 */
func PlusVideoWatch(vid int) int {
	if err := db.Table("videodata").Where("vid = ?", vid).UpdateColumn("watch", gorm.Expr("watch + ?", 1)).Error; err != nil {
		LogUtils.Logger("[数据库错误SQL Error] 在用户增加视频查看次数时产生：" + err.Error())
		return 0
	}
	return 1
}

func IsRelation(uid int, vid int) int {
	var relationData RelationModel.Relation
	db.Table("relation").Where("uid = ? AND vid = ?", uid, vid).Take(&relationData)
	return relationData.Id
}
