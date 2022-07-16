/*
 * @Descripttion: 转码库工具 - FFmpeg Utils
 * @Author: William Wu
 * @Date: 2022/06/16 上午 11:47
 * @LastEditors: William Wu
 * @LastEditTime: 2022/06/16 上午 11:47
 */

// 注意:此部分工具都必须要求系统环境内有FFmpeg的环境
// Note: this part of tools must have ffmpeg environment in the system environment

package FFmpegUtils

import "VideoHubGo/utils/SystemUtils"

/**
 * @Descripttion: 截取视频封面 - Capture video cover
 * @Author: William Wu
 * @Date: 2022/06/16 下午 03:39
 * @Param: filePath (string)
 * @Param: coverPath (string)
 * @Return: result (int)
 */
func GetVideoCover(filePath string, coverPath string) int {
	cmd := "ffmpeg -y -ss 30 -i " + filePath + " -f mjpeg -t 1 -s 520x290 " + coverPath + ".png"
	osSys := SystemUtils.GetSystemInfo()
	if osSys == 1 {
		_, err := SystemUtils.RunLinuxCommand(cmd)
		if err != nil {
			return 0
		}
	} else if osSys == 2 {
		_, err := SystemUtils.RunWindowsCommand(cmd)
		if err != nil {
			return 0
		}
	} else {
		return 0
	}
	return 1
}

/**
 * @Descripttion: 二次编码视频 - Recode Video
 * @Author: William Wu
 * @Date: 2022/06/16 下午 03:44
 * @Param: filePath (string)
 * @Param: outPath (string)
 * @Return: result (int)
 */
func RecodeVideo(filePath string, outPath string) int {
	cmd := "ffmpeg -i " + filePath + " -movflags faststart -acodec copy -vcodec copy " + outPath
	osSys := SystemUtils.GetSystemInfo()
	if osSys == 1 {
		_, err := SystemUtils.RunLinuxCommand(cmd)
		if err != nil {
			return 0
		}
	} else if osSys == 2 {
		_, err := SystemUtils.RunWindowsCommand(cmd)
		if err != nil {
			return 0
		}
	} else {
		return 0
	}
	return 1
}
