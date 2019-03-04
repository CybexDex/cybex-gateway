package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"coding.net/bobxuyang/cy-gateway-BN/app"
	"coding.net/bobxuyang/cy-gateway-BN/controllers/adminsrv"
	model "coding.net/bobxuyang/cy-gateway-BN/models"
	"coding.net/bobxuyang/cy-gateway-BN/utils"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	version   string
	githash   string
	buildtime string
	branch    string
)

func main() {
	v := flag.Bool("v", false, "version")
	flag.Parse()
	if *v {
		fmt.Printf("version: %s_%s_%s, build time: %s\n", version, branch, githash, buildtime)
		return
	}

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// init config
	utils.InitConfig()
	logDir := viper.GetString("adminsrv.log_dir")
	logLevel := viper.GetString("adminsrv.log_level")
	// init logger
	utils.InitLog(logDir, logLevel, "[admin]")
	utils.Infof("version: %s_%s_%s, build time: %s", version, branch, githash, buildtime)

	// init db
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetString("database.port")
	dbUser := viper.GetString("database.user")
	dbPassword := viper.GetString("database.pass")
	dbName := viper.GetString("database.name")
	model.InitDB(dbHost, dbPort, dbUser, dbPassword, dbName)

	// init route
	router := mux.NewRouter()
	router.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}).Methods("GET")

	router.HandleFunc("/api/balance/{account}", adminsrv.GetBalance).Methods("GET")

	router.HandleFunc("/api/account/new", adminsrv.CreateAccount).Methods("POST")
	router.HandleFunc("/api/account/login", adminsrv.Authenticate).Methods("POST")
	router.HandleFunc("/api/debug/info", adminsrv.DebugInfo).Methods("GET")

	router.HandleFunc("/api/jadepool/new", adminsrv.CreateJadepool).Methods("POST")
	router.HandleFunc("/api/jadepool/{id}", adminsrv.UpdateJadepool).Methods("PUT")
	router.HandleFunc("/api/jadepool/{id}", adminsrv.GetJadepool).Methods("GET")
	router.HandleFunc("/api/jadepool", adminsrv.GetAllJadepool).Methods("GET")
	router.HandleFunc("/api/jadepool/{id}", adminsrv.DeleteJadepool).Methods("DELETE")

	router.HandleFunc("/api/blockchain/new", adminsrv.CreateBlockchain).Methods("POST")
	router.HandleFunc("/api/blockchain/{id}", adminsrv.UpdateBlockchain).Methods("PUT")
	router.HandleFunc("/api/blockchain/{id}", adminsrv.GetBlockchain).Methods("GET")
	router.HandleFunc("/api/blockchain", adminsrv.GetAllBlockchain).Methods("GET")
	router.HandleFunc("/api/blockchain/{id}", adminsrv.DeleteBlockchain).Methods("DELETE")

	router.HandleFunc("/api/asset/new", adminsrv.CreateAsset).Methods("POST")
	router.HandleFunc("/api/asset/{id}", adminsrv.UpdateAsset).Methods("PUT")
	router.HandleFunc("/api/asset/{id}", adminsrv.GetAsset).Methods("GET")
	router.HandleFunc("/api/asset", adminsrv.GetAllAsset).Methods("GET")
	router.HandleFunc("/api/asset/{id}", adminsrv.DeleteAsset).Methods("DELETE")

	// init middleware
	router.Use(app.NewLoggingMiddle(utils.GetLogger()))
	//router.Use(app.JwtAuthentication) //attach JWT auth middleware
	corsHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Content-Length", "Authorization", "Accept", "X-Requested-With", "Current-Page"})
	corsMethods := handlers.AllowedMethods([]string{"GET", "DELETE", "PUT", "POST", "PATCH", "OPTIONS"})
	corsOrigins := handlers.AllowedOrigins([]string{"*"})
	handler := handlers.CORS(corsHeaders, corsMethods, corsOrigins)(router)

	listenAddr := viper.GetString("adminsrv.listen_addr")
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
