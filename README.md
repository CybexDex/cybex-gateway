# gateway

## 功能简介

​    连接瑶池和cybex链： 

​        充值: 瑶池-> gateway - cybex 上给用户转账 ->cybex转账 

​        提现: 用户转账给cybex上的网关账户->gateway 扫链获得信息 ->发送瑶池提现请求  

## 使用

本项目使用go mod. 

git clone 到 GOPATH之外 

设置 GO111MODULE=on

参考 config 下 template.yaml 编写自己的 xxx.yaml

​env=xxx go run cmd/all/main.go

默认使用 dev.yaml

## 配置

参见 config/template.yaml 中的注释

其中 eccPub 和 eccPri 是瑶池相关的，可以直接配置。也可以写成 seed__{eccpubkey} ,seed_{eccprikey} 。
{keyname}  表示 在seed中配置的keyname。 seed__ 前缀会使用 seed中的数据。

考虑点

1. 是直接连接瑶池，还是使用saas。如果使用saas，useSass选项设置为true，sassserver下的配置需要配置。否则jpserver下需要配置
2. userserver 是供客户端访问的api，是否要验签和支持跨域在 userserver 下配置
3. seed 用于存储敏感信息。
4. 使用不使用微信高级 wx.enable 设置为false

## 订单子阶段
充值
  jp,order,cyborder,done

提现
  order,jp,jp_sended,done

### 错误处理
充值

  cyborder 子状态 init processing pending failed
  其中failed可以安全重试。
  processing 且无 sig可以安全重试

提现

  jp 子状态 fail 可以安全重试。但是可能是瑶池钱不够，或者提币数额高于瑶池配置等重试也无效的原因。所以最好除了网络fail就人工排查下。
  