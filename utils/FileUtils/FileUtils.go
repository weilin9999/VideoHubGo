/*
 * @Descripttion: 文件操作工具 - File Utils
 * @Author: William Wu
 * @Date: 2022/06/17 上午 11:11
 * @LastEditors: William Wu
 * @LastEditTime: 2022/06/17 上午 11:11
 */

package FileUtils

import (
	"VideoHubGo/models/AdminModel"
	"VideoHubGo/utils/LogUtils"
	"io/ioutil"
	"os"
)

/**
 * @Descripttion: 移动文件到指定位置 - Move file to specified location
 * @Author: William Wu
 * @Date: 2022/06/17 上午 11:12
 * @Param: filePath (string)
 * @Param: specifyPath (string)
 * @Return: result (int)
 */
func MoveFile(oldFilePath string, specifyPath string) int {
	err := os.Rename(oldFilePath, specifyPath)
	if err != nil {
		LogUtils.Logger("[IO操作异常 - IO operation exception] 移动文件时出现异常 Move file operation exception：" + err.Error())
		return 0
	}
	return 1
}

/**
 * @Descripttion: 删除指定文件 - Delete file to specified location
 * @Author: William Wu
 * @Date: 2022/06/17 上午 11:12
 * @Param: filePath (string)
 * @Return: result (int)
 */
func DeleteFile(filePath string) int {
	err := os.Remove(filePath)
	if err != nil {
		LogUtils.Logger("[IO操作异常 - IO operation exception] 删除文件时出现异常 Delete file operation exception：" + err.Error())
		return 0
	}
	return 1
}

/**
 * @Descripttion: 输出文件夹下的所有文件 - Traverse All File
 * @Author: William Wu
 * @Date: 2022/06/17 上午 11:51
 * @Param: filePath (string)
 * @Return: result ([]AdminModel.AdminFileList)
 */
func TraverseFile(filePath string) ([]AdminModel.AdminFileList, error) {
	var fileArray []AdminModel.AdminFileList
	ioRead, err := ioutil.ReadDir(filePath)
	if err != nil {
		LogUtils.Logger("[IO操作异常 - IO operation exception] 遍历文件时出现异常 Traverse file operation exception：" + err.Error())
		return fileArray, err
	}
	for index, file := range ioRead {
		temp := AdminModel.AdminFileList{
			Id:   index,
			Name: file.Name(),
			Size: file.Size(),
			Date: file.ModTime(),
		}
		fileArray = append(fileArray, temp)
	}
	return fileArray, err
}
