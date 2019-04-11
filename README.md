# bbb-gateway

为了bbb的定制化网关

## 开始

go run cmd/bbb/main.go

会检查配置，报出错误。

## 需求

### 1. 充值

```
server:jpserver
关系:瑶池
config:
  瑶池地址
  端口   

other_config: 瑶池需要配置Jpsrv地址
检查:向瑶池发起一笔提现。
瑶池数据:coinName

支持币种
server:userserver
config:
  USDT
    充值:
      地址: seed__gatewayin
      password: "seed__pass:gatewayin"
      转化为:
        NB
        JADE.USDT:seed__maker
        CYB::1
    提现
      地址: gatewayout
      password: "seed__pass:gatewayin"
      等待提现
        JADE.USDT
        发送到 gatewayin
  cybex链:
    地址
other_config: seed 中 gatewayin gatewayout 密码
关系:seed
config:
  seed 地址
  seed lib库,从seed服务器获取
  commandkey
检查:用户名密码正确。

检查ecc
  base64 hex
worker
config:
  worker时间
```

2. 提现
3. 充提列表
4. 错误探测和处理

任务入口有 server 和 worker。 cmd会去调用server或者worker


## TODO

日志格式
请求日志格式
  请求返回记录
BN请求记录
ecc
jp_withdraw_eos
jp_deposit_eos
jp_deposit_eth

float 计算问题



