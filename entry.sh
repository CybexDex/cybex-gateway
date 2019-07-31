nohup ./bin/all >./log/gateway.log 2>&1 &
nohup ./bin/admin >./log/admin.log 2>&1 &
sleep 10
tail -f ./log/gateway.log
