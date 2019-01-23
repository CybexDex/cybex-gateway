package main

import (
	"net/http"
	"os"
	"time"

	"git.coding.net/bobxuyang/cy-gateway-BN/app"
	"git.coding.net/bobxuyang/cy-gateway-BN/controllers"
	"git.coding.net/bobxuyang/cy-gateway-BN/utils"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	githash   string
	buildtime string
	branch    string
)

func main() {
	// 配置初始化日志
	logDir := os.Getenv("log_dir")
	logLevel := os.Getenv("log_level")
	utils.InitLog(logDir, logLevel)
	utils.Infof("build info: %s_%s_%s", buildtime, branch, githash)

	// 配置路由
	router := mux.NewRouter()
	router.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}).Methods("GET")
	router.HandleFunc("/api/account/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/account/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/debug/info", controllers.DebugInfo).Methods("GET")

	router.HandleFunc("/api/blockchain/new", controllers.CreateBlockchain).Methods("POST")
	router.HandleFunc("/api/blockchain/{id}", controllers.UpdateBlockchain).Methods("PUT")
	router.HandleFunc("/api/blockchain/{id}", controllers.GetBlockchain).Methods("GET")
	router.HandleFunc("/api/blockchain", controllers.GetAllBlockchain).Methods("GET")
	router.HandleFunc("/api/blockchain/{id}", controllers.DeleteBlockchain).Methods("DELETE")

	// 配置中间件
	router.Use(app.JwtAuthentication) //attach JWT auth middleware
	router.Use(app.NewLoggingMiddle(utils.GetLogger()))

	listenAddr := os.Getenv("listen_addr")
	if len(listenAddr) == 0 {
		listenAddr = ":8082"
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
