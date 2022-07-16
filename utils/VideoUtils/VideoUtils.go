/*
 * @Descripttion: 视频信息工具 - Video Utils
 * @Author: William Wu
 * @Date: 2022/06/16 下午 01:30
 * @LastEditors: William Wu
 * @LastEditTime: 2022/06/16 下午 01:30
 */
package VideoUtils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
)

type entityBoxHeader struct {
	Size       uint32
	FourccType [4]byte
	Size64     uint64
}

/**
 * @Descripttion: 获取视频总时长 - Get Video Total Time
 * @Author: William Wu
 * @Date: 2022/06/16 下午 01:35
 * @Param: filePath (string)
 * @Return: time (string)
 */
func GetVideoTotalTime(filePath string) (duration string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "文件未找到 - File found for", err
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			return
		}
	}(file)

	var (
		info      = make([]byte, 0x10)
		boxHeader entityBoxHeader
		offset    int64 = 0
	)
	// 获取结构偏移
	for {
		_, err = file.ReadAt(info, offset)
		if err != nil {
			return
		}
		boxHeader = getHeaderBoxInfo(info)
		fourccType := getFourccType(boxHeader)
		if fourccType == "moov" {
			break
		}
		// 有一部分mp4 mdat尺寸过大需要特殊处理
		if fourccType == "mdat" {
			if boxHeader.Size == 1 {
				offset += int64(boxHeader.Size64)
				continue
			}
		}
		offset += int64(boxHeader.Size)
	}
	// 获取move结构开头一部分
	moveStartBytes := make([]byte, 0x100)
	_, err = file.ReadAt(moveStartBytes, offset)
	if err != nil {
		return
	}
	// 定义timeScale与Duration偏移
	timeScaleOffset := 0x1C
	durationOffset := 0x20
	timeScale := binary.BigEndian.Uint32(moveStartBytes[timeScaleOffset : timeScaleOffset+4])
	Duration := binary.BigEndian.Uint32(moveStartBytes[durationOffset : durationOffset+4])
	return resolveTime(Duration / timeScale), nil
}

/**
 * @Descripttion: 获取头信息 - Get Header Box Info
 * @Author: William Wu
 * @Date: 2022/06/16 下午 01:35
 * @Param: data ([]byte)
 * @Return: entityBoxHeader
 */
func getHeaderBoxInfo(data []byte) (boxHeader entityBoxHeader) {
	buf := bytes.NewBuffer(data)
	_ = binary.Read(buf, binary.BigEndian, &boxHeader)
	return
}

/**
 * @Descripttion: 获取信息头类型 - Get Fourcc Type
 * @Author: William Wu
 * @Date: 2022/06/16 下午 01:36
 * @Param: entityBoxHeader
 * @Return: fourccType (string)
 */
func getFourccType(boxHeader entityBoxHeader) (fourccType string) {
	fourccType = string(boxHeader.FourccType[:])
	return
}

/**
 * @Descripttion: 将秒转成时分秒格式 - Resolve Time
 * @Author: William Wu
 * @Date: 2022/06/16 下午 01:37
 * @Param: seconds (uint32)
 * @Return: time (string)
 */
func resolveTime(seconds uint32) string {
	var (
		h, m, s string
	)
	var day = seconds / (24 * 3600)
	hour := (seconds - day*3600*24) / 3600
	minute := (seconds - day*24*3600 - hour*3600) / 60
	second := seconds - day*24*3600 - hour*3600 - minute*60
	h = strconv.Itoa(int(hour))
	if hour < 10 {
		h = "0" + strconv.Itoa(int(hour))
	}
	m = strconv.Itoa(int(minute))
	if minute < 10 {
		m = "0" + strconv.Itoa(int(minute))
	}
	s = strconv.Itoa(int(second))
	if second < 10 {
		s = "0" + strconv.Itoa(int(second))
	}
	if h == "00" {
		return fmt.Sprintf("%s:%s", m, s)
	} else {
		return fmt.Sprintf("%s:%s:%s", h, m, s)
	}
}
