# gateway

## 功能简介

​    连接瑶池和cybex链： 

​        充值: 瑶池-> gateway - cybex 上给用户转账 ->cybex转账 

​        提现: 用户转账给cybex上的网关账户->gateway 扫链获得信息 ->发送瑶池提现请求  

## 使用

​    本项目使用go mod. 

​    git clone 到 GOPATH之外 

​	设置 GO111MODULE=on

​	参考 config 下 template.yaml 编写自己的 xxx.yaml

​	env=xxx go run cmd/bbb/main.go

## 配置

参见 config/template.yaml 中的注释

其中 eccPub 和 eccPri 是瑶池相关的，可以直接配置。也可以写成 seed__{eccpubkey} ,seed_{eccprikey}
{keyname}  表示 在seed中配置的keyname。 seed__ 前缀会使用 seed中的数据。

