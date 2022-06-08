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
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

var conn = RedisUtils.RedisPool.Get()

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
		conn.Send("ZADD", "videodata", v.Vid, string(jsonTemp))
	}
	conn.Flush()
	_, err := conn.Receive()
	if err != nil {
		LogUtils.Logger("[Redis操作]Redis存储失败-VideoWriteList操作中")
	}
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

	res2, err := redis.Values(conn.Do("zrangebyscore", "videodata", (count - endId), (count - startId)))
	var tempVideo []VideoModel.VideoRe

	for _, v := range res2 {
		tempdata := VideoModel.VideoRe{}
		errd := json.Unmarshal(v.([]byte), &tempdata)
		if errd != nil {
			fmt.Println(err)
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
	_, err := conn.Do("set", "videocount", count)
	if err != nil {
		LogUtils.Logger("[Redis操作] 存储视频总数时出现异常：" + err.Error())
	}
}

/**
 * @Descripttion: 从Redis获取视频总数 - Total number of videos obtained from redis
 * @Author: William Wu
 * @Date: 2022/06/07 下午 03:21
 */
func VideoGetCount() int {
	count, err := redis.Int(conn.Do("get", "videocount"))
	if err != nil {
		LogUtils.Logger("[Redis操作] 获取视频总数时出现异常：" + err.Error())
	}
	return count
}

/**
 * @Descripttion: 从已经存储的Redis视频数据获取集合总数 - Total number of collections obtained from stored redis video data
 * @Author: William Wu
 * @Date: 2022/06/07 下午 03:21
 */
func GetReidsVideoListCount() int {
	count, err := redis.Int(conn.Do("zcard", "videodata"))
	if err != nil {
		LogUtils.Logger("[Redis操作] 获取Rediso中的视频总数时出现异常：" + err.Error())
	}
	return count
}

/**
 * @Descripttion: 存入分类视频总数 - Save Class Video Count
 * @Author: William Wu
 * @Date: 2022/06/05 下午 01:14
 * @Param: count (int)
 */
func VideoSaveClassCountList(cid int, count int) {
	_, err := conn.Do("set", "classcount"+strconv.Itoa(cid), count)
	if err != nil {
		LogUtils.Logger("[Redis操作] 存储视频类型总数时出现异常：" + err.Error())
	}
}

/**
 * @Descripttion: 从Redis获取类型视频总数 - Total number of Class videos obtained from redis
 * @Author: William Wu
 * @Date: 2022/06/07 下午 03:21
 */
func VideoGetClassCount(cid int) int {
	count, err := redis.Int(conn.Do("get", "classcount"+strconv.Itoa(cid)))
	if err != nil {
		LogUtils.Logger("[Redis操作] 获取视频类型总数时出现异常：" + err.Error())
	}
	return count
}
