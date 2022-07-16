/*
 * @Descripttion: 视频类型模型 - Video Class Model
 * @Author: William Wu
 * @Date: 2022/05/27 下午 03:43
 * @LastEditors: William Wu
 * @LastEditTime: 2022/05/27 下午 03:43
 */
package ClassModel

import (
	"time"
)

/**
 * @Descripttion: 视频类型模型 - Video Class Model
 * @Author: William Wu
 * @Date: 2022/05/27 下午 03:45
 */
type Class struct {
	Cid         int    `gorm:"primaryKey" json:"cid"`
	Classname   string `json:"classname"`
	Isdelete    int    `json:"isdelete"`
	Create_Time time.Time
	Update_Time time.Time
}

/**
 * @Descripttion: 回传类型数据 - Return type data
 * @Author: William Wu
 * @Date: 2022/06/03 下午 02:31
 */
type ClassRe struct {
	Cid       int    `json:"cid"`
	Classname string `json:"classname"`
}

/**
 * @Descripttion: 后台视频分类模型 - Admin video class struct
 * @Author: William Wu
 * @Date: 2022/07/05 下午 09:17
 */
type ClassAdminSearch struct {
	Cid       int    `json:"cid"`
	Classname string `json:"classname"`
	Page      int    `json:"page"`
	Size      int    `json:"size"`
}

/**
 * @Descripttion: 后台分类模型 - Admin class struct
 * @Author: William Wu
 * @Date: 2022/07/05 下午 09:18
 */
type ClassAdminRe struct {
	Cid         int       `json:"cid"`
	Classname   string    `json:"classname"`
	Isdelete    int       `json:"isdelete"`
	Create_Time time.Time `json:"create_time"`
}
