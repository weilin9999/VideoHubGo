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
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

var conn = RedisUtils.RedisPool.Get()

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
		conn.Send("ZADD", "classdata", v.Cid, string(jsonTemp))
	}
	conn.Flush()
	_, err := conn.Receive()
	if err != nil {
		LogUtils.Logger("[Redis操作]Redis存储失败-ClassWriteCache操作中")
	}
}

/**
 * @Descripttion: Redis读取视频分类数据 - Read Video Class Data From Rredis
 * @Author: William Wu
 * @Date: 2022/06/03 下午 02:50
 * @Return: ClassModel ClassRe
 */
func ClassGetListCache() []ClassModel.ClassRe {

	res2, err := redis.Values(conn.Do("zcard", "classdata"))
	var tempClass []ClassModel.ClassRe

	for _, v := range res2 {
		tempdata := ClassModel.ClassRe{}
		errd := json.Unmarshal(v.([]byte), &tempdata)
		if errd != nil {
			fmt.Println(err)
		}
		tempClass = append(tempClass, tempdata)
	}

	return tempClass
}
