/*
 * @Descripttion: 图片模型 - Photo Model
 * @Author: William Wu
 * @Date: 2022/07/08 下午 10:48
 * @LastEditors: William Wu
 * @LastEditTime: 2022/07/08 下午 10:48
 */
package PhotoModel

import "time"

/**
 * @Descripttion: 图片基本模型 - Photo base model
 * @Author: William Wu
 * @Date: 2022/07/08 下午 11:03
 */
type PhotoModel struct {
	Pid         int    `gorm:"primaryKey" json:"pid"`
	Plast       string `json:"plast"`
	Isdelete    int    `json:"isdelete"`
	Create_Time time.Time
	Update_Time time.Time
}

/**
 * @Descripttion: 图片数据返回模型 - Photo data return model
 * @Author: William Wu
 * @Date: 2022/07/08 下午 11:03
 */
type PhotoModelRe struct {
	Pid         int       `json:"pid"`
	Plast       string    `json:"plast"`
	Isdelete    int       `json:"isdelete"`
	Create_Time time.Time `json:"create_time"`
}

/**
 * @Descripttion: 图片所搜模型 - Photo search model
 * @Author: William Wu
 * @Date: 2022/07/08 下午 11:03
 */
type PhotoSearch struct {
	Page int `json:"page"`
	Size int `json:"size"`
}
