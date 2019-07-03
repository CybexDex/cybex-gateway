## 增加新币种

在 coins 新建 xx.yaml

xx.yaml 的内容参考XTZ.yaml

coin=xx  go run cmd/db/add_asset.go

第一次执行会创建

后面执行如果name已经存在会更新资产

暂时不支持删除资产，但是可以disable掉