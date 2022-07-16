/*
 * @Descripttion: 图片缓存层 - Photo Cache
 * @Author: William Wu
 * @Date: 2022/07/08 下午 11:09
 * @LastEditors: William Wu
 * @LastEditTime: 2022/07/08 下午 11:09
 */
package PhotoCache

import (
	"VideoHubGo/models/PhotoModel"
	"VideoHubGo/utils/LogUtils"
	"VideoHubGo/utils/RedisUtils"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

var conn = RedisUtils.RedisClient
var ctx = context.Background()

/**
 * @Descripttion: 图片数据存入Redis - Photo Data Save Redis
 * @Author: William Wu
 * @Date: 2022/07/08 下午 11:10
 * @Param: PhotoModel PhotoModelRe
 */
func PhotoWriteListCache(photoData []PhotoModel.PhotoModelRe) {
	for k, v := range photoData {
		jsonTemp, err := json.Marshal(photoData[k])
		if err != nil {
			LogUtils.Logger("[Redis操作]Json转换失败")
		}
		conn.ZAdd(ctx, "photodata", &redis.Z{
			Score:  float64(v.Pid),
			Member: jsonTemp,
		})
	}
	conn.Persist(ctx, "photodata")
}

/**
 * @Descripttion: Redis读取图片数据 - Read Photo Data From Rredis
 * @Author: William Wu
 * @Date: 2022/07/08 下午 11:12
 * @Param: page (int)
 * @Param: size (int)
 * @Return: VideoModel VideoRe
 */
func PhotoGetListCache(page int, size int) []PhotoModel.PhotoModelRe {

	startId := (page * size) - size
	endId := page * size
	count := PhotoGetCount() + 1000

	res2, err := conn.ZRangeByScore(ctx, "photodata", &redis.ZRangeBy{
		Min: strconv.Itoa((count - endId)),
		Max: strconv.Itoa((count - startId - 1)),
	}).Result()
	var tempPhoto []PhotoModel.PhotoModelRe
	countArray := len(res2) - 1
	for i := countArray; i >= 0; i-- {
		v := res2[i]
		tempdata := PhotoModel.PhotoModelRe{}
		errd := json.Unmarshal([]byte(v), &tempdata)
		if errd != nil {
			LogUtils.Logger("[Redis操作] 获取图片数据时出现异常：" + err.Error())
		}
		tempPhoto = append(tempPhoto, tempdata)
	}
	return tempPhoto
}

/**
 * @Descripttion: 存入图片总数 - Save Photo Count
 * @Author: William Wu
 * @Date: 2022/07/08 下午 11:14
 * @Param: count (int)
 */
func PhotoSaveCountList(count int) {
	err := conn.Set(ctx, "photocount", count, time.Hour*2).Err()
	if err != nil {
		LogUtils.Logger("[Redis报错] 在存储图片总数时出错：" + err.Error())
	}
	conn.Persist(ctx, "photocount")
}

/**
 * @Descripttion: 从Redis获取图片总数 - Total number of photos obtained from redis
 * @Author: William Wu
 * @Date: 2022/07/08 下午 11:16
 */
func PhotoGetCount() int {
	count, err := conn.Get(ctx, "photocount").Result()
	if err != nil {
		return 0
	}
	tCount, _ := strconv.Atoi(count)
	return tCount
}

/**
 * @Descripttion: 从已经存储的Redis图片数据获取集合总数 - Total number of collections obtained from stored redis photo data
 * @Author: William Wu
 * @Date: 2022/07/08 下午 11:17
 */
func GetReidsPhotoListCount() int {
	count, err := conn.ZCard(ctx, "photodata").Result()
	if err != nil {
		return -1
	}
	return int(count)
}

/**
 * @Descripttion: 删除图片缓存 - Delete photo cache
 * @Author: William Wu
 * @Date: 2022/07/08 下午 11:17
 */
func PhotoDeleteCaches() {
	_, err := conn.Del(ctx, "photodata", "photocount").Result()
	if err != nil {
		LogUtils.Logger("[Redis操作]删除photodata,photocount时异常：" + err.Error())
	}
}
