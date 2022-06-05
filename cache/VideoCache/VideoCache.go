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
	res2, err := redis.Values(conn.Do("zrangebyscore", "videodata", (startId + 1000), (endId + 999)))
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
