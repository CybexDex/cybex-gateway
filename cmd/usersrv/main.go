package main

import (
	"net/http"
	"strings"
	"time"

	"git.coding.net/bobxuyang/cy-gateway-BN/app"
	"git.coding.net/bobxuyang/cy-gateway-BN/controllers/usersrv"
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
			// utils.Respond(w, utils.Message(false, "Invalid/Malformed auth token err:3"), http.StatusForbidden)
			// return
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
func main() {
	// 配置初始化日志
	utils.InitConfig()
	// init loggger
	logDir := viper.GetString("usersrv.log_dir")
	logLevel := viper.GetString("usersrv.log_level")
	utils.InitLog(logDir, logLevel)
	utils.Infof("build info: %s_%s_%s", buildtime, branch, githash)
	// route
	router := mux.NewRouter()
	router.Use(app.NewLoggingMiddle(utils.GetLogger()))
	authrouter := router.PathPrefix("/").Subrouter()
	authrouter.Use(authMiddleware)

	router.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}).Methods("GET")

	router.HandleFunc("/login", usersrv.Login).Methods("POST")
	authrouter.HandleFunc("/asset", usersrv.AllAsset).Methods("GET")
	authrouter.HandleFunc("/deposit_address/{user}/{asset}", usersrv.DepositAddress).Methods("GET")

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
