# make file

curDir := $(shell pwd)
date := $(shell date +%Y%m%d-%H:%M:%S)
githash := $(shell git log -1 --format="%h")
gitbranch := $(shell git rev-parse --abbrev-ref HEAD)

#######################################all server#####################################
.PHONY: buildAll buildAllLinux
buildAll: buildJPSrv buildAdminSrv buildCybSrv buildOrderSrv buildUserSrv

buildAllLinux: buildJPSrvLinux buildAdminSrvLinux buildCybSrvLinux buildOrderSrvLinux buildUserSrvLinux

scpAllDev: scpJPSrvDev scpAdminSrvDev scpCybSrvDev scpUserSrvDev scpOrderSrvDev devRestart
#######################################jpSrv#########################################
.PHONY: buildJPSrv buildJPSrvLinux startJPSrv scpJPSrvDev devRestart
buildJPSrv:
	@echo "build jpsrv......"
	@(cd ${curDir}/cmd/jpsrv;\
	go build -v -ldflags "-X main.githash=$(githash) -X main.buildtime=$(date) -X main.branch=$(gitbranch)" -o ${curDir}/bin/jpsrv)

buildJPSrvLinux:
	@echo "build jpsrv linux......"
	@(cd ${curDir}/cmd/jpsrv;\
	GOOS=linux GOARCH=amd64 go build -v -ldflags "-X main.githash=$(githash) -X main.buildtime=$(date) -X main.branch=$(gitbranch)" -o ${curDir}/bin/jpsrv_linux_amd64)

startJPSrv: buildJPSrv
	@echo "start jpsrv......"
	@(${curDir}/bin/jpsrv;)
	
scpJPSrvDev: buildJPSrvLinux
	@echo "scp jpsrv......"
	@(scp bin/jpsrv_linux_amd64 root@39.98.58.238:~/gateway/bin/jpsrv_)


#######################################adminSrv#####################################
.PHONY: buildAdminSrv buildAdminSrvLinux startAdminSrv scpAdminSrvDev
buildAdminSrv:
	@echo "build adminsrv......"
	@(cd ${curDir}/cmd/adminsrv;\
	go build -v -ldflags "-X main.githash=$(githash) -X main.buildtime=$(date) -X main.branch=$(gitbranch)" -o ${curDir}/bin/adminsrv)

buildAdminSrvLinux:
	@echo "build adminsrv linux......"
	@(cd ${curDir}/cmd/adminsrv;\
	GOOS=linux GOARCH=amd64 go build -v -ldflags "-X main.githash=$(githash) -X main.buildtime=$(date) -X main.branch=$(gitbranch)" -o ${curDir}/bin/adminsrv_linux_amd64)

startAdminSrv: buildAdminSrv
	@echo "start adminsrv......"
	@(${curDir}/bin/adminsrv;)

scpAdminSrvDev: buildAdminSrvLinux
	@echo "scp adminsrv......"
	@(scp bin/adminsrv_linux_amd64 root@39.98.58.238:~/gateway/bin/adminsrv_)
devRestart: 
	@(ssh root@39.98.58.238 "cd /root/gateway && sh /root/gateway/startAll.sh")
  
#######################################cybSrv#########################################
.PHONY: buildCybSrv startCybSrv buildCybSrvLinux scpCybSrvDev
buildCybSrv:
	@echo "build cybsrv......"
	@(cd ${curDir}/cmd/cybsrv;\
	go build -v -ldflags "-X main.githash=$(githash) -X main.buildtime=$(date) -X main.branch=$(gitbranch)" -o ${curDir}/bin/cybsrv)

buildCybSrvLinux:
	@echo "build cybsrv linux......"
	@(cd ${curDir}/cmd/cybsrv;\
	GOOS=linux GOARCH=amd64 go build -v -ldflags "-X main.githash=$(githash) -X main.buildtime=$(date) -X main.branch=$(gitbranch)" -o ${curDir}/bin/cybsrv_linux_amd64)

startCybSrv: buildCybSrv
	@echo "start cybsrv......"
	@(${curDir}/bin/cybsrv;)

scpCybSrvDev: buildCybSrvLinux
	@echo "scp cybsrv......"
	@(scp bin/cybsrv_linux_amd64 root@39.98.58.238:~/gateway/bin/cybsrv_)

#######################################orderSrv#########################################
.PHONY: buildOrderSrv startOrderSrv buildOrderSrvLinux scpOrderSrvDev
buildOrderSrv:
	@echo "build ordersrv......"
	@(cd ${curDir}/cmd/ordersrv;\
	go build -v -ldflags "-X main.githash=$(githash) -X main.buildtime=$(date) -X main.branch=$(gitbranch)" -o ${curDir}/bin/ordersrv)

buildOrderSrvLinux:
	@echo "build ordersrv linux......"
	@(cd ${curDir}/cmd/ordersrv;\
	GOOS=linux GOARCH=amd64 go build -v -ldflags "-X main.githash=$(githash) -X main.buildtime=$(date) -X main.branch=$(gitbranch)" -o ${curDir}/bin/ordersrv_linux_amd64)

startOrderSrv: buildOrderSrv
	@echo "start ordersrv......"
	@(${curDir}/bin/ordersrv;)

scpOrderSrvDev: buildOrderSrvLinux
	@echo "scp ordersrv......"
	@(scp bin/ordersrv_linux_amd64 root@39.98.58.238:~/gateway/bin/ordersrv_)

#######################################userSrv#########################################
.PHONY: buildUserSrv startUserSrv buildUserSrvLinux scpUserSrvDev
buildUserSrv:
	@echo "build usersrv......"
	@(cd ${curDir}/cmd/usersrv;\
	go build -v -ldflags "-X main.githash=$(githash) -X main.buildtime=$(date) -X main.branch=$(gitbranch)" -o ${curDir}/bin/usersrv)

buildUserSrvLinux:
	@echo "build usersrv linux......"
	@(cd ${curDir}/cmd/usersrv;\
	GOOS=linux GOARCH=amd64 go build -v -ldflags "-X main.githash=$(githash) -X main.buildtime=$(date) -X main.branch=$(gitbranch)" -o ${curDir}/bin/usersrv_linux_amd64)

startUserSrv: buildUserSrv
	@echo "start usersrv......"
	@(${curDir}/bin/usersrv;)

scpUserSrvDev: buildUserSrvLinux
	@echo "scp usersrv......"
	@(scp bin/usersrv_linux_amd64 root@39.98.58.238:~/gateway/bin/usersrv_)

#######################################swagger#########################################
scpSwagger:
	@echo "scp swagger......"
	@(rsync -avzc --delete doc/swagger/ root@39.98.58.238:~/swagger/static/swagger/)