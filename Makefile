# make file

curDir := $(shell pwd)
date := $(shell date +%Y%m%d-%H:%M:%S)
githash := $(shell git log -1 --format="%h")
gitbranch := $(shell git rev-parse --abbrev-ref HEAD)

#######################################buildall#####################################
.PHONY: buildAll buildAllLinux
buildAll: buildJPSrv buildAdminSrv

buildAllLinux: buildJPSrvLinux buildAdminSrvLinux

#######################################jpSrv#########################################
.PHONY: buildJPSrv buildJPSrvLinux startJPSrv scpJPSrvDev
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
	@(scp bin/jpsrv_linux_amd64 root@39.98.58.238:~/jpsrv/jpsrv_)


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
	@(scp bin/adminsrv_linux_amd64 root@39.98.58.238:~/adminsrv/adminsrv_)

#######################################cybSrv#########################################
.PHONY: buildCybSrv startCybSrv
buildCybSrv:
	@echo "build jpsrv......"
	@(cd ${curDir}/cmd/cybsrv;\
	go build -v -ldflags "-X main.githash=$(githash) -X main.buildtime=$(date) -X main.branch=$(gitbranch)" -o ${curDir}/bin/cybsrv)

startCybSrv: buildCybSrv
	@echo "start cybsrv......"
	@(${curDir}/bin/cybsrv;)

#######################################orderSrv#########################################
.PHONY: buildOrderSrv startOrderSrv
buildOrderSrv:
	@echo "build orderSrv......"
	@(cd ${curDir}/cmd/orderSrv;\
	go build -v -ldflags "-X main.githash=$(githash) -X main.buildtime=$(date) -X main.branch=$(gitbranch)" -o ${curDir}/bin/orderSrv)

startOrderSrv: buildCybSrv
	@echo "start orderSrv......"
	@(${curDir}/bin/orderSrv;)

#######################################userSrv#########################################
.PHONY: buildUserSrv startUserSrv
buildUserSrv:
	@echo "build userSrv......"
	@(cd ${curDir}/cmd/userSrv;\
	go build -v -ldflags "-X main.githash=$(githash) -X main.buildtime=$(date) -X main.branch=$(gitbranch)" -o ${curDir}/bin/userSrv)

startUserSrv: buildCybSrv
	@echo "start userSrv......"
	@(${curDir}/bin/userSrv;)