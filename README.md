# VideoHubGo
The projectis  make GoLang program,this project uses GO+Gin+Gorm+Redis 该项目使用Go语言编写，使用GO+Gin+Gorm+Redis框架

## 项目详细 - Object detail
本项目采用GO语言来编写后端，其中使用了Gorm、Gin、Redis的技术来实现后端的基本的数据库、缓存与请求的需求，其中Redis主要运用于视频数据的缓存上和分类数据的缓存上，高频的数据是具有Redis的缓存与高效存储，数据库主要使用MySQL进行数据的存储，视频转码技术运用了FFmpeg的工具，代码运行已经经过测试与优化。

## 对接项目前端
前端所使用了VUE3+Vite的方案，具体的前端请查看我的仓库有对应独立的后台和用户端（注：用户端视频上传功能并未完善请等待，后台已全部完善）

## 搭建详细
项目需要有FFmpeg的环境，MySQL 5.6及以上，Redis服务，数据库文件在本项目的SQL文件夹中
