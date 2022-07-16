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
	"time"
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
	db.Select("videodata.vid,videodata.uid,videodata.detail,videodata.watch,videodata.vtime,videodata.create_time,userdata.username").Table("videodata").Joins("JOIN userdata ON videodata.uid = userdata.uid").Where("videodata.isdelete = ?", 0).Order("videodata.vid DESC").Limit(size).Offset(offset).Find(&videoData)
	return videoData
}

/**
 * @Descripttion: 查询视频数据总数 - Count Video List
 * @Author: William Wu
 * @Date: 2022/06/05 下午 01:03
 * @Return: count (int)
 */
func GetCountVideoList() int {
	var count int64
	db.Select("vid").Table("videodata").Where("isdelete = ?", 0).Count(&count)
	return int(count)
}

/**
 * @Descripttion: 从数据中查找分类的视频数据 - Find classified video data from data
 * @Author: William Wu
 * @Date: 2022/06/05 下午 06:28
 * @Param: cid (int)
 * @Param: size (int)
 * @Param: offset (int)
 * @Return: VideoModel VideoRe
 */
func FindVideoInClass(cid int, size int, offset int) []VideoModel.VideoRe {
	var videoData []VideoModel.VideoRe
	db.Select("videodata.vid,videodata.uid,videodata.detail,videodata.watch,videodata.vtime,videodata.create_time,userdata.username").Table("videodata").Joins("JOIN userdata ON videodata.uid = userdata.uid").Where("videodata.isdelete = ? AND videodata.cid = ?", 0, cid).Order("videodata.vid DESC").Limit(size).Offset(offset).Find(&videoData)
	return videoData
}

/**
 * @Descripttion: 查询相关类型视频数据总数 - Count Class Video List
 * @Author: William Wu
 * @Date: 2022/06/05 下午 01:03
 * @Param: cid (int)
 * @Return: count (int)
 */
func GetCountVideoClassList(cid int) int {
	var count int64
	db.Select("vid").Table("videodata").Where("isdelete = ? and cid = ?", 0, cid).Count(&count)
	return int(count)
}

/**
 * @Descripttion: 搜索视频内容 - Search Video List
 * @Author: William Wu
 * @Date: 2022/06/08 下午 03:58
 * @Param: cid (int)
 * @Param: key (string)
 * @Param: size (int)
 * @Param: offset (int)
 * @Return: VideoModel VideoRe
 */
func SearchVideoList_Class(cid int, key string, size int, offset int) ([]VideoModel.VideoRe, int) {
	var videoData []VideoModel.VideoRe
	var count int64
	if cid == 0 {
		db.Select("videodata.vid,videodata.uid,videodata.detail,videodata.watch,videodata.vtime,videodata.create_time,userdata.username").Table("videodata").Joins("JOIN userdata ON videodata.uid = userdata.uid").Where("videodata.isdelete = ? AND videodata.detail LIKE ?", 0, key+"%").Order("videodata.vid DESC").Limit(size).Offset(offset).Count(&count).Find(&videoData)
	} else {
		db.Select("videodata.vid,videodata.uid,videodata.detail,videodata.watch,videodata.vtime,videodata.create_time,userdata.username").Table("videodata").Joins("JOIN userdata ON videodata.uid = userdata.uid").Where("videodata.isdelete = ? AND videodata.cid = ? AND videodata.detail LIKE ?", 0, cid, "%"+key+"%").Order("videodata.vid DESC").Limit(size).Offset(offset).Count(&count).Find(&videoData)
	}
	return videoData, int(count)
}

/**
 * @Descripttion: 视频预上传存入数据库 - Video pre upload stored in database
 * @Author: William Wu
 * @Date: 2022/06/15 下午 05:22
 * @Param: detail (string)
 * @Return: vid (int)
 */
func UploadVideo(uid int, detail string, cid int, vtime string) int {
	videoData := VideoModel.Video{}
	videoData.Uid = uid
	videoData.Detail = detail
	videoData.Cid = cid
	videoData.Vtime = vtime
	videoData.Create_Time = time.Now()
	videoData.Update_Time = time.Now()
	delVid := FindIsdeleteVideoVid()
	if delVid != 0 {
		videoData.Vid = delVid
		videoData.Isdelete = 0
		videoData.Watch = 0
		db.Table("videodata").Save(&videoData)
	} else {
		db.Table("videodata").Create(&videoData)
	}
	return videoData.Vid
}

/**
 * @Descripttion: 查询已经软删除的视频VID - Query Isdeleted Video VID
 * @Author: William Wu
 * @Date: 2022/06/15 下午 05:32
 * @Return: vid (int)
 */
func FindIsdeleteVideoVid() int {
	videoData := VideoModel.Video{}
	db.Table("videodata").Where("isdelete = ?", 1).First(&videoData)
	return videoData.Vid
}
