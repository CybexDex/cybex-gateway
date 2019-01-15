#! /bin/sh

BIN_NAME=$2
CUR_DIR=`pwd`
BIN_NAME_LINUX=$BIN_NAME\_linux_amd64

usage()
{
    echo $0 dev/build/buildLinux
}

build()
{
    DATE=`date +%Y%m%d-%H:%M:%S`
    GITHASH=`git log -1 --format="%h"`
    GITBRANCH=`git rev-parse --abbrev-ref HEAD`

    BINPATH=$1
    cd $CUR_DIR/cmd/$BIN_NAME
    go build -v -ldflags "-X main.githash=$GITHASH -X main.buildtime=$DATE -X main.branch=$GITBRANCH" -o $BINPATH 
    cd -
}

start()
{
    cd $CUR_DIR/config/$BIN_NAME
    $1
}

case $1 in
    "build")
        BINPATH=$CUR_DIR/bin/$BIN_NAME
        build $BINPATH;;
    "dev")
        BINPATH=$CUR_DIR/bin/$BIN_NAME
        build $BINPATH
        start $BINPATH;;
    "buildLinux")
        BINPATH=$CUR_DIR/bin/$BIN_NAME_LINUX
        GOOS=linux GOARCH=amd64 build $BINPATH;;
    *)
        usage;;
esac

exit 0