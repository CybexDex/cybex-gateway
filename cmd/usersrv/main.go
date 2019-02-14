package main

import (
	"net/http"
	"time"

	"git.coding.net/bobxuyang/cy-gateway-BN/app"
	"git.coding.net/bobxuyang/cy-gateway-BN/controllers/usersrvc"
	"git.coding.net/bobxuyang/cy-gateway-BN/utils"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

var (
	githash   string
	buildtime string
	branch    string
)

func main() {
	// 配置初始化日志
	utils.InitConfig()
	// init loggger
	logDir := viper.GetString("usersrv.log_dir")
	logLevel := viper.GetString("usersrv.log_level")
	utils.InitLog(logDir, logLevel)
	utils.Infof("build info: %s_%s_%s", buildtime, branch, githash)
	// 配置路由
	router := mux.NewRouter()
	router.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}).Methods("GET")
	router.HandleFunc("/login", usersrvc.Login).Methods("POST")
	router.HandleFunc("/asset", usersrvc.AllAsset).Methods("GET")
	router.HandleFunc("/deposit_address/{user}/{asset}", usersrvc.DepositAddress).Methods("GET")
	router.Use(app.NewLoggingMiddle(utils.GetLogger()))

	listenAddr := viper.GetString("usersrv.listen_addr")
	utils.Infof("%s", listenAddr)
	if len(listenAddr) == 0 {
		listenAddr = ":8081"
	}
	server := &http.Server{
		Addr:         listenAddr,
		Handler:      router,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	server.SetKeepAlivesEnabled(false)

	// 启动server
	utils.Infof("listen on: %s", listenAddr)
	gracehttp.Serve(server)
}
