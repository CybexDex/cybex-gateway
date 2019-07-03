nohup ./bin/all >./log/gateway.log 2>&1 &
sleep 10
tail -f ./log/gateway.log
