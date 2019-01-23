curDir := $(shell pwd)
date := $(shell date +%Y%m%d-%H:%M:%S)
githash := $(shell git log -1 --format="%h")
gitbranch := $(shell git rev-parse --abbrev-ref HEAD)

#######################################buildall#####################################
.PHONY: buildAll buildAllLinux
buildAll: buildJPSrv buildAdminSrv

buildAllLinux: buildJPSrvLinux buildAdminSrvLinux

#######################################jpSrv#####################################
.PHONY: buildJPSrv buildJPSrvLinux startJPSrv 
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
	@(cd ${curDir}/config/jpsrv; \
	${curDir}/bin/jpsrv;)

#######################################adminSrv#####################################
.PHONY: buildAdminSrv buildAdminSrvLinux startAdminSrv
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
	@(cd ${curDir}/config/adminsrv; \
	${curDir}/bin/adminsrv;)