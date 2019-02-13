package main

import (
	"net/http"
	"os"
	"time"

	"git.coding.net/bobxuyang/cy-gateway-BN/app"
	"git.coding.net/bobxuyang/cy-gateway-BN/controllers"
	"git.coding.net/bobxuyang/cy-gateway-BN/utils"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/handlers"
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
	logDir := viper.GetString("adminsrv.log_dir")
	logLevel := viper.GetString("adminsrv.log_level")
	// init logger
	utils.InitLog(logDir, logLevel)
	utils.Infof("build info: %s_%s_%s", buildtime, branch, githash)

	// init route
	router := mux.NewRouter()
	router.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}).Methods("GET")
	router.HandleFunc("/api/account/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/account/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/debug/info", controllers.DebugInfo).Methods("GET")

	router.HandleFunc("/api/jadepool/new", controllers.CreateJadepool).Methods("POST")
	router.HandleFunc("/api/jadepool/{id}", controllers.UpdateJadepool).Methods("PUT")
	router.HandleFunc("/api/jadepool/{id}", controllers.GetJadepool).Methods("GET")
	router.HandleFunc("/api/jadepool", controllers.GetAllJadepool).Methods("GET")
	router.HandleFunc("/api/jadepool/{id}", controllers.DeleteJadepool).Methods("DELETE")

	router.HandleFunc("/api/blockchain/new", controllers.CreateBlockchain).Methods("POST")
	router.HandleFunc("/api/blockchain/{id}", controllers.UpdateBlockchain).Methods("PUT")
	router.HandleFunc("/api/blockchain/{id}", controllers.GetBlockchain).Methods("GET")
	router.HandleFunc("/api/blockchain", controllers.GetAllBlockchain).Methods("GET")
	router.HandleFunc("/api/blockchain/{id}", controllers.DeleteBlockchain).Methods("DELETE")

	router.HandleFunc("/api/asset/new", controllers.CreateAsset).Methods("POST")
	router.HandleFunc("/api/asset/{id}", controllers.UpdateAsset).Methods("PUT")
	router.HandleFunc("/api/asset/{id}", controllers.GetAsset).Methods("GET")
	router.HandleFunc("/api/asset", controllers.GetAllAsset).Methods("GET")
	router.HandleFunc("/api/asset/{id}", controllers.DeleteAsset).Methods("DELETE")

	// init middleware
	router.Use(app.NewLoggingMiddle(utils.GetLogger()))
	//router.Use(app.JwtAuthentication) //attach JWT auth middleware
	corsHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Content-Length", "Authorization", "Accept", "X-Requested-With", "Current-Page"})
	corsMethods := handlers.AllowedMethods([]string{"GET", "PUT", "POST", "PATCH", "OPTIONS"})
	corsOrigins := handlers.AllowedOrigins([]string{"*"})
	handler := handlers.CORS(corsHeaders, corsMethods, corsOrigins)(router)

	listenAddr := os.Getenv("adminsrv.listen_addr")
	if len(listenAddr) == 0 {
		listenAddr = ":8082"
	}
	server := &http.Server{
		Addr:         listenAddr,
		Handler:      handler,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	server.SetKeepAlivesEnabled(false)

	// start server
	utils.Infof("listen on: %s", listenAddr)
	gracehttp.Serve(server)
}
