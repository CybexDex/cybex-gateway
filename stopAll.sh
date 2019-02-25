#!/bin/bash

# exit if error
set -e

# stop jpsrv
echo "stop jpsrv..."
ps aux | grep -ie "bin/jpsrv" | grep -v 'grep'| awk '{print $2}' | xargs kill -9

# stop adminsrv
echo "stop adminsrv..."
ps aux | grep -ie "bin/adminsrv" | grep -v 'grep'| awk '{print $2}' | xargs kill -9

# stop cybsrv
echo "stop cybsrv..."
ps aux | grep -ie "bin/cybsrv" | grep -v 'grep'| awk '{print $2}' | xargs kill -9

# stop ordersrv
echo "stop ordersrv..."
ps aux | grep -ie "bin/ordersrv" | grep -v 'grep'| awk '{print $2}' | xargs kill -9

# stop usersrv
echo "stop usersrv..."
ps aux | grep -ie "bin/usersrv" | grep -v 'grep'| awk '{print $2}' | xargs kill -9