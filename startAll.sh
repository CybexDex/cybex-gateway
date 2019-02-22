#!/bin/bash
set -e

if [ ! -d "./log" ]; then 
    mkdir log
fi

make buildJPSrv buildAdminSrv buildCybSrv buildOrderSrv buildUserSrv

echo "start jpsrv..."
ps aux | grep -ie "bin/jpsrv" | grep -v 'grep'| awk '{print $2}' | xargs kill -9
nohup ./bin/jpsrv >./log/jpsrv.log 2>&1 &

echo "start adminsrv..."
ps aux | grep -ie "bin/adminsrv" | grep -v 'grep'| awk '{print $2}' | xargs kill -9
nohup ./bin/adminsrv >./log/adminsrv.log 2>&1 &

echo "start cybsrv..."
ps aux | grep -ie "bin/cybsrv" | grep -v 'grep'| awk '{print $2}' | xargs kill -9
nohup ./bin/cybsrv >./log/cybsrv.log 2>&1 &

echo "start ordersrv..."
ps aux | grep -ie "bin/ordersrv" | grep -v 'grep'| awk '{print $2}' | xargs kill -9
nohup ./bin/ordersrv >./log/ordersrv.log 2>&1 &

echo "start usersrv..."
ps aux | grep -ie "bin/usersrv" | grep -v 'grep'| awk '{print $2}' | xargs kill -9
nohup ./bin/usersrv >./log/usersrv.log 2>&1 &