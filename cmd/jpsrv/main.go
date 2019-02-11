package main

import (
	"net/http"
	"time"

	"git.coding.net/bobxuyang/cy-gateway-BN/app"
	"git.coding.net/bobxuyang/cy-gateway-BN/controllers"
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
	// init config
	utils.InitConfig()
	// init loggger
	logDir := viper.GetString("jpsrv.log_dir")
	logLevel := viper.GetString("jpsrv.log_level")
	utils.InitLog(logDir, logLevel)
	utils.Infof("build info: %s_%s_%s", buildtime, branch, githash)

	// init route
	router := mux.NewRouter()
	router.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}).Methods("GET")
	router.HandleFunc("/api/order/noti", controllers.NotiOrder).Methods("POST")
	router.HandleFunc("/api/address/new", controllers.GetNewAddress).Methods("GET")
	router.HandleFunc("/api/order/send", controllers.SendOrder).Methods("POST")
	router.Use(app.NewLoggingMiddle(utils.GetLogger()))

	listenAddr := viper.GetString("jpsrv.listen_addr")
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

	// start server
	utils.Infof("listen on: %s", listenAddr)
	gracehttp.Serve(server)
}
