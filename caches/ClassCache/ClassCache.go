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
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
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
			LogUtils.Logger("[Redis操作]Json转换失败")
		}
		conn.ZAdd(ctx, "classdata", &redis.Z{
			Score:  float64(v.Cid),
			Member: jsonTemp,
		})
	}
	conn.ExpireAt(ctx, "classdata", time.Now().Add(2*time.Hour))
}

/**
 * @Descripttion: Redis读取视频分类数据 - Read Video Class Data From Rredis
 * @Author: William Wu
 * @Date: 2022/06/03 下午 02:50
 * @Return: ClassModel ClassRe
 */
func ClassGetListCache() []ClassModel.ClassRe {
	res2, err := conn.ZRange(ctx, "classdata", 0, -1).Result()
	var tempClass []ClassModel.ClassRe
	for _, v := range res2 {
		tempdata := ClassModel.ClassRe{}
		errd := json.Unmarshal([]byte(v), &tempdata)
		if errd != nil {
			fmt.Println(err)
		}
		tempClass = append(tempClass, tempdata)
	}

	return tempClass
}
