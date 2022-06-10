/*
 * @Descripttion: 用户关系模型 - Relation Model
 * @Author: William Wu
 * @Date: 2022/05/27 下午 03:47
 * @LastEditors: William Wu
 * @LastEditTime: 2022/05/27 下午 03:47
 */
package RelationModel

/**
 * @Descripttion: 用户关系模型 - Relation Model
 * @Author: William Wu
 * @Date: 2022/05/27 下午 03:48
 */
type Relation struct {
	Id       int `gorm:"primaryKey" json:"id"`
	Uid      int `json:"uid"`
	Vid      int `json:"vid"`
	Isdelete int
}

/**
 * @Descripttion: Relation普通类型请求体 - Relation common type request struct
 * @Author: William Wu
 * @Date: 2022/06/09 下午 03:25
 */
type RelationRequest struct {
	Page int `json:"page"`
	Size int `json:"size"`
}
