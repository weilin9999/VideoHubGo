# VideoHubGo
该项目使用Go语言编写，使用GO+Gin+Gorm+Redis框架， The projectis make GoLang program,this project uses GO+Gin+Gorm+Redis

## 项目详细
本项目采用GO语言来编写后端，其中使用了Gorm、Gin、Redis的技术来实现后端的基本的数据库、缓存与请求的需求，其中Redis主要运用于视频数据的缓存上和分类数据的缓存上，高频的数据是具有Redis的缓存与高效存储，数据库主要使用MySQL进行数据的存储，视频转码技术运用了FFmpeg的工具，代码运行已经经过测试与优化。

## 对接项目前端
前端所使用了VUE3+Vite的方案，具体的前端请查看我的仓库有对应独立的后台和用户端（注：用户端视频上传功能并未完善请等待，后台已全部完善）
Android端虽然已经开发完毕但是仅闭源不开源。

## 搭建详细
项目需要有FFmpeg的环境（主要是针对视频文件的处理），MySQL 5.6及以上（后续可能会转Oracle），Redis服务，数据库文件在本项目的SQL文件夹中
如果你想要Nginx进行转发后台的流量可进行Nginx的反向代理

## 项目详细功能
项目主要采用GO语言为主，前后端分离，前端包括有（后台、用户端），移动端仅开发了Android，
（目前后台功能已经全部完善，如有问题或者意见可反馈，用户端还为完善因写的太匆忙用户端仅仅只是展示）

### 1、路由
项目的路由文件都在router/Router.go文件中，里面包括了用户端的路由、后台的路由以及移动端的路由，
每条路由都有对应的注释以及对应的Controller，静态的资源、上传文件的规则以及Token的校验和JWT的中间件都在
Router内配置，Token和静态资源的配置都在configs/config.yaml文件中可修改。

### 2、控制器
项目的控制器全都在controllers的包内，每一个Controller都有一个对应的Services和Models以及Router，
Controller进行逻辑的处理反馈接收的数据。

### 3、业务层
项目的业务层都在services的包内，进行业务的处理，每一个service对应着一个controller以及models，
业务层是进行业务的逻辑和数据处理，可复用性高。

### 4、缓冲层
项目采用Redis与Mysql进行数据的读写分离，在该项目内caches的包内是处理Redis逻辑业务的处理，
缓冲热门数据以及业务的频繁读取的数据，从而达到业务的优化和提升用户的使用体验。

### 5、数据模型
项目的所有的结构体都在models的包内，每一个model都对应着一个controller和service。

### 6、配置文件
项目的配置文件都在configs/config.yaml中，内部有注释每一条每个小项的具体对应的配置是什么，
具体的详细信息请到改文件内查看即可。

## 最后
项目的雏形基本已经全部完善，如果有BUG请发issues，我会及时查看和解决。