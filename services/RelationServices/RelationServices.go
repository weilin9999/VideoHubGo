/*
 * @Descripttion: 用户收藏服务层 - Relation Services
 * @Author: William Wu
 * @Date: 2022/06/09 下午 03:27
 * @LastEditors: William Wu
 * @LastEditTime: 2022/06/09 下午 03:27
 */
package RelationServices

import (
	"VideoHubGo/models/RelationModel"
	"VideoHubGo/models/VideoModel"
	"VideoHubGo/utils/DataBaseUtils"
	"VideoHubGo/utils/LogUtils"
)

var db = DataBaseUtils.GoDB()

/**
 * @Descripttion: 查询用户收藏信息集合 - Query user favorite information collection
 * @Author: William Wu
 * @Date: 2022/06/09 下午 03:29
 * @Param: uid (int)
 * @Param: size (int)
 * @Param: offset (int)
 * @Return: VideoModel Video, count
 */
func FindRelationByVideoList(uid int, size int, offset int) ([]VideoModel.VideoRe, int) {
	var videoData []VideoModel.VideoRe
	var count int64
	db.Select("videodata.vid,videodata.detail,videodata.watch,videodata.vtime,videodata.cid,videodata.create_time").Table("relation").Joins("LEFT JOIN videodata ON relation.vid = videodata.vid").Where(" relation.uid = ? AND relation.isdelete = ? AND videodata.isdelete = ?", uid, 0, 0).Order("videodata.vid DESC").Limit(size).Offset(offset).Count(&count).Find(&videoData)
	return videoData, int(count)
}

/**
 * @Descripttion: 查询用户收藏分类信息集合- Query user favorite information class collection
 * @Author: William Wu
 * @Date: 2022/06/10 上午 10:01
 * @Param: uid (int)
 * @Param: cid (int)
 * @Param: size (int)
 * @Param: offset (int)
 * @Return: VideoModel Video, count
 */
func FindRelationByVideoList_Class(uid int, cid int, size int, offset int) ([]VideoModel.VideoRe, int) {
	var videoData []VideoModel.VideoRe
	var count int64
	db.Select("videodata.vid,videodata.detail,videodata.watch,videodata.vtime,videodata.cid,videodata.create_time").Table("relation").Joins("LEFT JOIN videodata ON relation.vid = videodata.vid").Where(" relation.uid = ? AND videodata.cid = ? AND relation.isdelete = ? AND videodata.isdelete = ?", uid, cid, 0, 0).Order("videodata.vid DESC").Limit(size).Offset(offset).Count(&count).Find(&videoData)
	return videoData, int(count)
}

/**
 * @Descripttion: 搜索用户收藏信息集合(包含分类) - Search user favorite information collection (including classification)
 * @Author: William Wu
 * @Date: 2022/06/10 下午 03:35
 * @Param: uid (int)
 * @Param: cid (int)
 * @Param: key (string)
 * @Param: size (int)
 * @Param: offset (int)
 * @Return: VideoModel Video, count
 */
func SearchRelationByVideoList_Class(uid int, cid int, key string, size int, offset int) ([]VideoModel.VideoRe, int) {
	var videoData []VideoModel.VideoRe
	var count int64
	if cid == 0 {
		db.Select("videodata.vid,videodata.detail,videodata.watch,videodata.vtime,videodata.cid,videodata.create_time").Table("relation").Joins("LEFT JOIN videodata ON relation.vid = videodata.vid").Where(" relation.uid = ? AND relation.isdelete = ? AND videodata.isdelete = ? AND videodata.detail LIKE ?", uid, 0, 0, key+"%").Order("videodata.vid DESC").Limit(size).Offset(offset).Count(&count).Find(&videoData)
	} else {
		db.Select("videodata.vid,videodata.detail,videodata.watch,videodata.vtime,videodata.cid,videodata.create_time").Table("relation").Joins("LEFT JOIN videodata ON relation.vid = videodata.vid").Where(" relation.uid = ? AND videodata.cid = ? AND relation.isdelete = ? AND videodata.isdelete = ? AND videodata.detail LIKE ?", uid, cid, 0, 0, "%"+key+"%").Order("videodata.vid DESC").Limit(size).Offset(offset).Count(&count).Find(&videoData)
	}
	return videoData, int(count)
}

/**
 * @Descripttion: 添加用户收藏 - Add Relation
 * @Author: William Wu
 * @Date: 2022/06/10 下午 03:38
 * @Param: uid (int)
 * @Param: vid (int)
 * @Return: success (int)
 */
func AddRelation(uid int, vid int) int {
	var reData RelationModel.Relation
	reData.Uid = uid
	reData.Vid = vid
	delId := IsDeleteRelation()
	if delId != 0 {
		reData.Id = delId
		reData.Isdelete = 0
		if err := db.Table("relation").Save(&reData).Error; err != nil {
			LogUtils.Logger("[数据库错误 SQL Error]在操作用户添加收藏更新时产生异常：" + err.Error())
			return 0
		}
	} else {
		if err := db.Table("relation").Create(&reData).Error; err != nil {
			LogUtils.Logger("[数据库错误 SQL Error]在操作用户添加收藏时产生异常：" + err.Error())
			return 0
		}
	}
	return 1
}

/**
 * @Descripttion: 取消用户收藏 - Delete Relation
 * @Author: William Wu
 * @Date: 2022/06/10 下午 03:39
 * @Param: uid (int)
 * @Param: vid (int)
 * @Return: success (int)
 */
func DeleteRelation(uid int, vid int) int {
	if err := db.Table("relation").Where("uid = ? AND vid = ?", uid, vid).Update("isdelete", 1).Error; err != nil {
		LogUtils.Logger("[数据库错误 SQL Error]在操作用户添加收藏时产生异常：" + err.Error())
		return 0
	}
	return 1
}

/**
 * @Descripttion: 查询已被软删除的收藏信息 - Query IsDelete Relation
 * @Author: William Wu
 * @Date: 2022/06/10 下午 03:59
 * @Return: id (int)
 */
func IsDeleteRelation() int {
	relationData := RelationModel.Relation{}
	db.Table("relation").Where("isdelete = ?", 1).First(&relationData)
	return relationData.Id
}

/**
 * @Descripttion: 查询该用户是否已收藏 - Query User Is Relation
 * @Author: William Wu
 * @Date: 2022/06/10 下午 04:09
 * @Param: uid (int)
 * @Param: vid (int)
 * @Return: id (int)
 */
func IsRelation(uid int, vid int) int {
	relationData := RelationModel.Relation{}
	db.Table("relation").Where("uid = ? AND vid = ? AND isdelete = ?", uid, vid, 0).First(&relationData)
	return relationData.Id
}
