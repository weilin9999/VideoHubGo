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
	Id       int `gorm:"primaryKey"`
	Uid      int
	Vid      int
	Isdelete int
}
