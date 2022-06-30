/*
 * @Descripttion: 视频缓存区 - Video Cache Area
 * @Author: William Wu
 * @Date: 2022/05/29 下午 05:04
 * @LastEditors: William Wu
 * @LastEditTime: 2022/05/29 下午 05:04
 */
package VideoCache

import (
	"VideoHubGo/models/VideoModel"
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
 * @Descripttion: 视频数据存入Redis - Video Data Save Redis
 * @Author: William Wu
 * @Date: 2022/05/29 下午 08:48
 * @Param: VideoModel Video
 */
func VideoWriteListCache(videoData []VideoModel.VideoRe) {
	for k, v := range videoData {
		jsonTemp, err := json.Marshal(videoData[k])
		if err != nil {
			LogUtils.Logger("[Redis操作]Json转换失败")
		}
		conn.ZAdd(ctx, "videodata", &redis.Z{
			Score:  float64(v.Vid),
			Member: jsonTemp,
		})
	}
	conn.ExpireAt(ctx, "videodata", time.Now().Add(2*time.Hour))
}

/**
 * @Descripttion: Redis读取视频数据 - Read Video Data From Rredis
 * @Author: William Wu
 * @Date: 2022/05/29 下午 08:48
 * @Param: page (int)
 * @Param: size (int)
 * @Return: VideoModel VideoRe
 */

func VideoGetListCache(page int, size int) []VideoModel.VideoRe {

	startId := (page * size) - size
	endId := page * size
	count := VideoGetCount() + 1000

	res2, err := conn.ZRangeByScore(ctx, "videodata", &redis.ZRangeBy{
		Min: strconv.Itoa((count - endId)),
		Max: strconv.Itoa((count - startId - 1)),
	}).Result()
	var tempVideo []VideoModel.VideoRe

	for _, v := range res2 {
		tempdata := VideoModel.VideoRe{}
		errd := json.Unmarshal([]byte(v), &tempdata)
		if errd != nil {
			LogUtils.Logger("[Redis操作] 获取视频数据时出现异常：" + err.Error())
		}
		tempVideo = append(tempVideo, tempdata)
	}
	return tempVideo
}

/**
 * @Descripttion: 存入视频总数 - Save Video Count
 * @Author: William Wu
 * @Date: 2022/06/05 下午 01:14
 * @Param: count (int)
 */
func VideoSaveCountList(count int) {
	err := conn.Set(ctx, "videocount", count, time.Hour*2).Err()
	if err != nil {
		LogUtils.Logger("[Redis报错] 在存储视频总数时出错：" + err.Error())
	}
	conn.ExpireAt(ctx, "videocount", time.Now().Add(2*time.Hour))
}

/**
 * @Descripttion: 从Redis获取视频总数 - Total number of videos obtained from redis
 * @Author: William Wu
 * @Date: 2022/06/07 下午 03:21
 */
func VideoGetCount() int {
	count, err := conn.Get(ctx, "videocount").Result()
	if err != nil {
		return 0
	}
	tCount, _ := strconv.Atoi(count)
	return tCount
}

/**
 * @Descripttion: 从已经存储的Redis视频数据获取集合总数 - Total number of collections obtained from stored redis video data
 * @Author: William Wu
 * @Date: 2022/06/07 下午 03:21
 */
func GetReidsVideoListCount() int {
	count, err := conn.ZCard(ctx, "videodata").Result()
	if err != nil {
		return -1
	}
	return int(count)
}

/**
 * @Descripttion: 存入分类视频总数 - Save Class Video Count
 * @Author: William Wu
 * @Date: 2022/06/05 下午 01:14
 * @Param: count (int)
 */
func VideoSaveClassCountList(cid int, count int) {
	err := conn.Set(ctx, "classcount"+strconv.Itoa(cid), count, time.Hour*2).Err()
	if err != nil {
		LogUtils.Logger("[Redis报错] 在存储视频总数时出错：" + err.Error())
	}

}

/**
 * @Descripttion: 从Redis获取类型视频总数 - Total number of Class videos obtained from redis
 * @Author: William Wu
 * @Date: 2022/06/07 下午 03:21
 */
func VideoGetClassCount(cid int) int {
	count, err := conn.Get(ctx, "classcount"+strconv.Itoa(cid)).Result()
	if err != nil {
		return 0
	}
	tCount, _ := strconv.Atoi(count)
	return tCount
}
