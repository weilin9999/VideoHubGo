/*
 * @Descripttion: 视频模型 - Video Model
 * @Author: William Wu
 * @Date: 2022/05/27 下午 03:33
 * @LastEditors: William Wu
 * @LastEditTime: 2022/05/27 下午 03:33
 */
package VideoModel

import (
	"time"
)

/**
 * @Descripttion: 视频模型 - Video Model
 * @Author: William Wu
 * @Date: 2022/05/27 下午 03:34
 */
type Video struct {
	Vid         int `gorm:"primaryKey" json:"vid"`
	Uid         int
	Detail      string
	Watch       int
	Vtime       string
	Cid         int
	Isdelete    int
	Create_Time time.Time
	Update_Time time.Time
}

/**
 * @Descripttion: 视频Redis模型 - Video Redis Model
 * @Author: William Wu
 * @Date: 2022/05/27 下午 03:34
 */
type VideoRe struct {
	Vid         int       `json:"vid"`
	Uid         int       `json:"uid"`
	Username    string    `json:"username"`
	Detail      string    `json:"detail"`
	Watch       int       `json:"watch"`
	Vtime       string    `json:"vtime"`
	Cid         int       `json:"cid"`
	Create_Time time.Time `json:"create_time"`
}

/**
 * @Descripttion: 视频数据请求模型 - Video Data Request Model
 * @Author: William Wu
 * @Date: 2022/05/27 下午 03:34
 */
type VideoRequest struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

/**
 * @Descripttion: 视频数据请求模型Class - Video Data Request Model Class
 * @Author: William Wu
 * @Date: 2022/05/27 下午 03:34
 */
type VideoRequestClass struct {
	Cid  int `json:"cid"`
	Page int `json:"page"`
	Size int `json:"size"`
}

/**
 * @Descripttion: 视频数据请求搜索模型Class - Video Data Request Search Model Class
 * @Author: William Wu
 * @Date: 2022/06/08 下午 03:53
 */
type VideoRequestSearch struct {
	Cid  int    `json:"cid"`
	Key  string `json:"key"`
	Page int    `json:"page"`
	Size int    `json:"size"`
}

/**
 * @Descripttion: 视频编辑模型 - Video edit struct
 * @Author: William Wu
 * @Date: 2022/06/30 下午 08:29
 */
type VideoEdit struct {
	Vid    int    `json:"vid"`
	Uid    int    `json:"uid"`
	Detail string `json:"detail"`
	Watch  int    `json:"watch"`
	Vtime  string `json:"vtime"`
	Cid    int    `json:"cid"`
}

/**
 * @Descripttion: 后台查询视频模型 - Admin search video struct
 * @Author: William Wu
 * @Date: 2022/07/05 下午 12:16
 */
type VideoAdminSearch struct {
	Cid    int    `json:"cid"`
	Vid    string `json:"vid"`
	Detail string `json:"detail"`
	Page   int    `json:"page"`
	Size   int    `json:"size"`
}

/**
 * @Descripttion: 后台视频数据模型 - Admin video struct
 * @Author: William Wu
 * @Date: 2022/07/05 下午 03:13
 */
type VideoAdminRe struct {
	Vid         int       `json:"vid"`
	Detail      string    `json:"detail"`
	Watch       int       `json:"watch"`
	Vtime       string    `json:"vtime"`
	Cid         int       `json:"cid"`
	Isdelete    int       `json:"isdelete"`
	Create_Time time.Time `json:"create_time"`
}
