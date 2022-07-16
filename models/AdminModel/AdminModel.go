/*
 * @Descripttion: 后台模型 - Admin Model
 * @Author: William Wu
 * @Date: 2022/07/05 下午 07:25
 * @LastEditors: William Wu
 * @LastEditTime: 2022/07/05 下午 07:25
 */
package AdminModel

import "time"

/**
 * @Descripttion: 视频组合模型 - Video group struct
 * @Author: William Wu
 * @Date: 2022/07/05 下午 11:08
 */
type VideoDelete struct {
	Vid int `json:"vid"`
}

/**
 * @Descripttion: 视频删除模型 - Video delete struct
 * @Author: William Wu
 * @Date: 2022/07/05 下午 11:09
 */
type VideoDeleteGroup struct {
	Group []VideoDelete `json:"group"`
}

/**
 * @Descripttion: 视频组合修改模型 - Video group edit struct
 * @Author: William Wu
 * @Date: 2022/07/05 下午 11:09
 */
type VideoBatchCid struct {
	Cid   int           `json:"cid"`
	Group []VideoDelete `json:"group"`
}

/**
 * @Descripttion: 分类删除模型 - Class delete struct
 * @Author: William Wu
 * @Date: 2022/07/05 下午 11:09
 */
type ClassDelete struct {
	Cid int `json:"cid"`
}

/**
 * @Descripttion: 分类组合删除模型 - Class delete group struct
 * @Author: William Wu
 * @Date: 2022/07/05 下午 11:10
 */
type ClassDeleteGroup struct {
	Group []ClassDelete `json:"group"`
}

/**
 * @Descripttion: 用户组合模型 - User group struct
 * @Author: William Wu
 * @Date: 2022/07/05 下午 11:11
 */
type UserDelete struct {
	Uid int `json:"uid"`
}

/**
 * @Descripttion: 用户删除组合模型 - User delete group struct
 * @Author: William Wu
 * @Date: 2022/07/05 下午 11:12
 */
type UserDeleteGroup struct {
	Group []UserDelete `json:"group"`
}

/**
 * @Descripttion: 用户管理员权限模型 - User admin authority Struct
 * @Author: William Wu
 * @Date: 2022/07/05 下午 11:21
 */
type UserIsadmin struct {
	Uid     int `json:"uid"`
	Isadmin int `json:"isadmin"`
}

/**
 * @Descripttion: 用户上传者权限模型 = User uploader authority Struct
 * @Author: William Wu
 * @Date: 2022/07/05 下午 11:22
 */
type UserIsuploader struct {
	Uid        int `json:"uid"`
	Isuploader int `json:"isuploader"`
}

/**
 * @Descripttion: 图片组合模型 - Photo group struct
 * @Author: William Wu
 * @Date: 2022/07/08 下午 11:53
 */
type PhotoDelete struct {
	Pid int `json:"pid"`
}

/**
 * @Descripttion: 图片删除组合模型 - Photo delete group struct
 * @Author: William Wu
 * @Date: 2022/07/08 下午 11:53
 */
type PhotoDeleteGroup struct {
	Group []PhotoDelete `json:"group"`
}

/**
 * @Descripttion: 后台数据模型 - Back-end data struct
 * @Author: William Wu
 * @Date: 2022/07/06 下午 06:43
 */
type AdminDashboard struct {
	Id      int `json:"id"`
	Re_Code int `json:"re_code"`
}

/**
 * @Descripttion: 后台文件列表 - Admin file list
 * @Author: William Wu
 * @Date: 2022/07/07 下午 01:08
 */
type AdminFileList struct {
	Id   int       `json:"id"`
	Name string    `json:"name"`
	Size int64     `json:"size"`
	Date time.Time `json:"date"`
}

/**
 * @Descripttion: 删除文件模型 - Delete struct
 * @Author: William Wu
 * @Date: 2022/07/07 下午 10:20
 */
type DeleteStruct struct {
	Name string `json:"name"`
}

/**
 * @Descripttion: 后台登录日志模型 - Admin login log struct
 * @Author: William Wu
 * @Date: 2022/07/16 上午 09:28
 */
type UserLog struct {
	Id          int       `gorm:"primaryKey" json:"id"`
	Uid         int       `json:"uid"`
	Account     string    `json:"account"`
	Username    string    `json:"username"`
	Ip          string    `json:"ip"`
	Isdelete    int       `json:"isdelete"`
	Create_Time time.Time `json:"create_time"`
	Update_Time time.Time `json:"update_time"`
}
