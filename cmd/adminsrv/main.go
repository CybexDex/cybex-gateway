package main

import (
	"net/http"
	"os"
	"time"

	"git.coding.net/bobxuyang/cy-gateway-BN/app"
	"git.coding.net/bobxuyang/cy-gateway-BN/utils"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var logger *utils.Logger

func main() {
	// 配置初始化日志
	logDir := os.Getenv("log_dir")
	logLevel := os.Getenv("log_level")
	utils.InitLog(logDir, logLevel)
	logger = utils.GetLogger()

	// 配置路由
	router := mux.NewRouter()
	router.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}).Methods("GET")

	// 配置中间件
	router.Use(app.JwtAuthentication) //attach JWT auth middleware
	router.Use(app.NewLoggingMiddle(logger))

	listenAddr := os.Getenv("listen_addr")
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
	logger.Infof("listen on: %s", listenAddr)
	gracehttp.Serve(server)
}
