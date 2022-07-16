/*
 * @Descripttion: 系统模型 - System Model
 * @Author: William Wu
 * @Date: 2022/06/17 下午 06:19
 * @LastEditors: William Wu
 * @LastEditTime: 2022/06/17 下午 06:19
 */
package SystemModel

/**
 * @Descripttion: 磁盘信息模型 - Disk MOdel
 * @Author: William Wu
 * @Date: 2022/06/17 下午 06:20
 */
type UsageStat struct {
	Path              string  `json:"path"`
	Fstype            string  `json:"fstype"`
	Total             uint64  `json:"total"`
	Free              uint64  `json:"free"`
	Used              uint64  `json:"used"`
	UsedPercent       float64 `json:"usedPercent"`
	InodesTotal       uint64  `json:"inodesTotal"`
	InodesUsed        uint64  `json:"inodesUsed"`
	InodesFree        uint64  `json:"inodesFree"`
	InodesUsedPercent float64 `json:"inodesUsedPercent"`
}
