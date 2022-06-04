/*
 * @Descripttion: 视频类型服务层 - Class Services
 * @Author: William Wu
 * @Date: 2022/06/03 下午 02:35
 * @LastEditors: William Wu
 * @LastEditTime: 2022/06/03 下午 02:35
 */
package ClassServices

import (
	"VideoHubGo/models/ClassModel"
	"VideoHubGo/utils/DataBaseUtils"
)

var db = DataBaseUtils.GoDB()

/**
 * @Descripttion: 查找所有的分类 - Find All Class
 * @Author: William Wu
 * @Date: 2022/06/03 下午 02:36
 * @Return: []ClassModel.ClassRe
 */
func FindAllClass() []ClassModel.ClassRe {
	var classData []ClassModel.ClassRe
	db.Select("cid,classname").Table("vclass").Find(&classData, "isdelete = ? ", 0)
	return classData
}
