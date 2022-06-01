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
	Cid         int
	Classname   string
	Isdelete    int
	Create_Time time.Time
	Update_Time time.Time
}
