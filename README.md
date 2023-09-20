# 自如/链家新房源实时提醒机器人🤖️

> 利用自如/链家网页版查询房源，解析HTML并分析房源，找到新上架房源并推送至钉钉或飞书群。

1. 首次初始化加载所选链接的所有房源（不通知）
2. 等待下次任务调度，调度周期时间由`taskInterval`控制
3. 执行任务，拿到最新房源数据，与上次房源集合进行比对
4. 已存在房源pass，新房源通知钉钉或飞书群

## 第一步：通过自如/链家获取房源检索地址

* 链家房源检索地址：https://bj.lianjia.com/zufang/
* 自如房源检索地址：https://www.ziroom.com/z/

## 第二步：运行GO程序

> 需自行编译（参考Golang交叉编译）。使用命令可参考：robot --help，查看提示信息。


```shell script
(base) ┌─[uzdz@uzdz] - [~/work/golang/rooms-inspect-robot] - [Tue Jul 27, 11:10]
└─[$] <git:(master*)> go run main.go --help                                                                                                                                                          ─╯
usage: main [<flags>] [<url>...]

Flags:
      --help                 Show context-sensitive help (also try --help-long and --help-man).
  -p, --notice="ding"        消息通知平台：ding（钉钉）、fs（飞书）
  -u, --noticeUrl=NOTICEURL  消息通知接口地址
  -k, --noticeKey="Home"     消息通知授权KEY（白名单）
  -t, --taskInterval=300     任务周期间隔时长，单位：秒
  -d, --proxyUrl=""          HTTP代理服务器配置，如果为空则不开启

Args:
  [<url>]  自如或链家网页版房源请求地址，支持录入多地址，多个地址通过`空格`分隔，复杂地址请进行UrlEncode操作后录入
```

以下进行举例：

> ./robot --notice=ding --noticeUrl='https://oapi.dingtalk.com/robot/send?access_token=xxx' --noticeKey=xxx https%3A%2F%2Fwww.ziroom.com%2Fz%2Fz2-s100011-r0%2F%3Fp%3Dx1%7C14%26cp%3D3000TO5000%26isOpen%3D1 https%3A%2F%2Fbj.lianjia.com%2Fditiezufang%2Fli651%2Fie1su1rt200600000001rp4%2F%3FshowMore%3D1

* UrlEncode工具网站：http://www.jsons.cn/urlencode/

## 第三步：钉钉通知

![](images/FCEF686C-A8A1-4FD5-AE75-038CA48A13E0.png)

# Go编译不同的平台文件

Golang 支持在一个平台下生成另一个平台可执行程序的交叉编译功能。

#### Mac下编译Linux, Windows平台的64位可执行程序：

* `Linux：`CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
* `Windows：`CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go

#### Linux下编译Mac, Windows平台的64位可执行程序：

* `Mac：`CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go
* `Windows：`CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go

#### Windows下编译Mac, Linux平台的64位可执行程序：

##### `Mac`

1. SET CGO_ENABLED=0
2. SET GOOS=darwin3
3. SET GOARCH=amd64
4. go build main.go

##### `Linux`

1. SET CGO_ENABLED=0
2. SET GOOS=linux
3. SET GOARCH=amd64
4. go build main.go

# License

This project is licensed under the [Apache v2.0 License](https://github.com/apache/skywalking-cli/blob/master/LICENSE).

# 免责声明

此软件程序用于替代人工耗时的检索房源过程，请勿修改代码中的网站保护策略。知法懂法，请参考[破坏计算机信息系统罪](https://www.66law.cn/zuiming/276.aspx)。
