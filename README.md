# bbb-gateway

为了bbb的定制化网关

## 开始

go run cmd/bbb/main.go

会检查配置，报出错误。

## 需求

1. 充值
2. 提现
3. 充提列表
4. 错误探测和处理

任务入口有 server 和 worker。 cmd会去调用server或者worker
