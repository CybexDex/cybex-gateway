package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"coding.net/bobxuyang/cy-gateway-BN/app"
	"coding.net/bobxuyang/cy-gateway-BN/controllers/usersrv"
	rep "coding.net/bobxuyang/cy-gateway-BN/help/singleton"
	model "coding.net/bobxuyang/cy-gateway-BN/models"
	"coding.net/bobxuyang/cy-gateway-BN/utils"
	"github.com/facebookgo/grace/gracehttp"
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

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			utils.Respond(w, utils.Message(false, "Missing auth token"), http.StatusForbidden)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			utils.Errorf("no token")
			utils.Respond(w, utils.Message(false, "Invalid/Malformed auth token err:1"), http.StatusForbidden)
			return
		}

		tokenPart := splitted[1]
		// check is tokenpart in db
		ok, err := usersrv.IsTokenOK(tokenPart)
		if err != nil {
			utils.Errorf("token err:%v", err)
			utils.Respond(w, utils.Message(false, "Invalid/Malformed auth token err:2"), http.StatusForbidden)
			return
		}
		if !ok {
			// TODO uncomment this
			utils.Respond(w, utils.Message(false, "Invalid/Malformed auth token err:3"), http.StatusForbidden)
			return
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
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

	// 配置初始化日志
	utils.InitConfig()
	// init loggger
	logDir := viper.GetString("usersrv.log_dir")
	logLevel := viper.GetString("usersrv.log_level")
	utils.InitLog(logDir, logLevel)
	utils.Infof("version: %s_%s_%s, build time: %s", version, branch, githash, buildtime)

	// init db
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetString("database.port")
	dbUser := viper.GetString("database.user")
	dbPassword := viper.GetString("database.pass")
	dbName := viper.GetString("database.name")
	model.InitDB(dbHost, dbPort, dbUser, dbPassword, dbName)
	rep.Init()

	// route
	router := mux.NewRouter()
	router.Use(app.NewLoggingMiddle(utils.GetLogger()))
	authrouter := router.PathPrefix("/").Subrouter()
	// authrouter.Use(authMiddleware)

	router.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}).Methods("GET")

	router.HandleFunc("/login", usersrv.Login).Methods("POST")
	// the supported assets
	authrouter.HandleFunc("/asset", usersrv.AllAsset).Methods("GET")
	// get deposit address
	authrouter.HandleFunc("/deposit_address/{user}/{asset}", usersrv.DepositAddress).Methods("GET")
	// get new deposit address
	authrouter.HandleFunc("/new_deposit_address/{user}/{asset}", usersrv.NewDepositAddress).Methods("GET")
	// verify_address
	authrouter.HandleFunc("/verify_address/{asset}/{address}", usersrv.VerifyAddress).Methods("GET")
	// deposit record
	authrouter.HandleFunc("/records/{user}", usersrv.Records).Methods("GET")
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
