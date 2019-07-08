version := 0.1.0
curDir := $(shell pwd)
date := $(shell date +%Y%m%d-%H:%M:%S)
githash := $(shell git log -1 --format="%h")
gitbranch := $(shell git rev-parse --abbrev-ref HEAD)
####
toDev: buildAllLinux scpAllDev scpDevConfig
##################
.PHONY: buildAll buildAllLinux scpDevConfig scpAllDev buildAdmin buildAdminLinux
buildAll:
	@echo "build all......"
	@(cd ${curDir}/cmd/all;\
	go build -v -ldflags "-X main.version=$(version) -X main.githash=$(githash) -X main.buildtime=$(date) -X main.branch=$(gitbranch)" -o ${curDir}/bin/all)
buildAllLinux:
	@echo "build all linux......"
	@(cd ${curDir}/cmd/all;\
	GOOS=linux GOARCH=amd64 go build -v -ldflags "-X main.version=$(version) -X main.githash=$(githash) -X main.buildtime=$(date) -X main.branch=$(gitbranch)" -o ${curDir}/bin/all)
buildAdmin:
	@echo "build admin......"
	@(cd ${curDir}/cmd/admin;\
	go build -v -ldflags "-X main.version=$(version) -X main.githash=$(githash) -X main.buildtime=$(date) -X main.branch=$(gitbranch)" -o ${curDir}/bin/admin)
buildAdminLinux:
	@echo "build admin linux......"
	@(cd ${curDir}/cmd/admin;\
	GOOS=linux GOARCH=amd64 go build -v -ldflags "-X main.version=$(version) -X main.githash=$(githash) -X main.buildtime=$(date) -X main.branch=$(gitbranch)" -o ${curDir}/bin/admin)
scpAllDev:
	@echo "scp All to dev......"
	#@(ssh root@39.98.58.238 "mv ~/bbb/config/dev.yaml ~/bbb/config/dev.yaml.bak")
	@(scp bin/all root@39.98.58.238:~/bbb/bin/)
scpDevConfig:
	@echo "scp dev.yaml......"
	@(ssh root@39.98.58.238 "mv ~/bbb/config/dev.yaml ~/bbb/config/dev.yaml.bak")
	@(scp config/dev.yaml root@39.98.58.238:~/bbb/config/dev.yaml)