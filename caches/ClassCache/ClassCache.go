/*
 * @Descripttion: 视频类型缓存区 - Class Cache
 * @Author: William Wu
 * @Date: 2022/06/03 下午 02:39
 * @LastEditors: William Wu
 * @LastEditTime: 2022/06/03 下午 02:39
 */
package ClassCache

import (
	"VideoHubGo/models/ClassModel"
	"VideoHubGo/utils/LogUtils"
	"VideoHubGo/utils/RedisUtils"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
)

var conn = RedisUtils.RedisClient
var ctx = context.Background()

/**
 * @Descripttion: 视频分类存入到Redis - Video Class Save Redis
 * @Author: William Wu
 * @Date: 2022/06/03 下午 02:41
 * @Param: ClassModel.ClassRe
 */
func ClassWriteCache(classData []ClassModel.ClassRe) {

	for k, v := range classData {
		jsonTemp, err := json.Marshal(classData[k])
		if err != nil {
			LogUtils.Logger("[Redis操作]Json转换失败-(ClassWriteCache)：" + err.Error())
		}
		conn.ZAdd(ctx, "classdata", &redis.Z{
			Score:  float64(v.Cid),
			Member: jsonTemp,
		})
	}
	conn.Persist(ctx, "classdata")
}

/**
 * @Descripttion: Redis读取视频分类数据 - Read Video Class Data From Rredis
 * @Author: William Wu
 * @Date: 2022/06/03 下午 02:50
 * @Return: ClassModel ClassRe
 */
func ClassGetListCache() []ClassModel.ClassRe {
	res2, err := conn.ZRange(ctx, "classdata", 0, -1).Result()
	if err != nil {
		LogUtils.Logger("[Redis操作]获取classdata时异常：" + err.Error())
	}
	var tempClass []ClassModel.ClassRe
	countArray := len(res2) - 1
	for i := countArray; i >= 0; i-- {
		v := res2[i]
		tempdata := ClassModel.ClassRe{}
		errd := json.Unmarshal([]byte(v), &tempdata)
		if errd != nil {
			LogUtils.Logger("[Redis操作]Json转换失败-(ClassGetListCache)：" + err.Error())
		}
		tempClass = append(tempClass, tempdata)
	}

	return tempClass
}

/**
 * @Descripttion: 删除分类缓存 - Delete class cache
 * @Author: William Wu
 * @Date: 2022/07/05 下午 10:31
 */
func ClassDeleteCaches() {
	_, err := conn.Del(ctx, "classdata").Result()
	if err != nil {
		LogUtils.Logger("[Redis操作]删除classdata时异常：" + err.Error())
	}
}
