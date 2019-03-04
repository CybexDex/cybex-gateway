package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"

	"coding.net/bobxuyang/cy-gateway-BN/models"
	"coding.net/bobxuyang/cy-gateway-BN/repository/jporder"

	"coding.net/bobxuyang/cy-gateway-BN/app"
	"coding.net/bobxuyang/cy-gateway-BN/controllers/jpsrv"
	"coding.net/bobxuyang/cy-gateway-BN/utils"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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
	// init loggger
	logDir := viper.GetString("jpsrv.log_dir")
	logLevel := viper.GetString("jpsrv.log_level")
	utils.InitLog(logDir, logLevel, "[jpsrv]")
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
		w.Write([]byte(`{"status": true}`))
	}).Methods("GET")
	router.HandleFunc("/api/order/noti", jpsrv.NotiOrder).Methods("POST")
	router.HandleFunc("/api/address/new", jpsrv.RequestNewAddress).Methods("GET")
	router.HandleFunc("/api/address/verify", jpsrv.VerifyBlockchainAddress).Methods("GET")
	router.HandleFunc("/api/confirmations", jpsrv.RequestConfirmations).Methods("GET")
	router.HandleFunc("/api/order/send", jpsrv.SendOrder).Methods("POST")
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

	go startHandleJPOrder()

	// start server
	utils.Infof("listen on: %s", listenAddr)
	gracehttp.Serve(server)
}

func startHandleJPOrder() {
	defer func() {
		if r := recover(); r != nil {
			utils.Errorf("r: %v, stack: %s", r, string(debug.Stack()))
		}
	}()

	jporderRepo := jporder.NewRepo(model.GetDB())
	for {
		for {
			jporder, err := jporderRepo.HoldingOne()
			if err != nil && err != gorm.ErrRecordNotFound {
				utils.Errorf("jporderRepo.HoldingOne error: %v", err)
				break
			}
			if jporder == nil || jporder.ID == 0 {
				break
			}

			utils.Infof("sending jporder(%d)", jporder.ID)
			data, err := jpsrv.DoSendOrder(jporder.ID)
			if err != nil {
				utils.Errorf("send jporder(%d) error: %v", jporder.ID, err)
				err = jporder.UpdateColumns(&model.JPOrder{
					Status: model.JPOrderStatusInit,
				})
				if err != nil {
					utils.Errorf("UpdateColumns jporder(%d) error: %v", jporder.ID, err)
				}
				time.Sleep(time.Second * 1)
				continue
			}

			resultBytes, err := json.Marshal(data)
			if err != nil {
				utils.Errorf("error: %v", err)
				err = jporder.UpdateColumns(&model.JPOrder{
					Status: model.JPOrderStatusInit,
				})
				if err != nil {
					utils.Errorf("UpdateColumns jporder(%d) error: %v", jporder.ID, err)
				}
				time.Sleep(time.Second * 1)
				continue
			}
			result := jpsrv.OrderNotiResult{}
			err = json.Unmarshal(resultBytes, &result)
			if err != nil {
				utils.Errorf("error: %v", err)
				err = jporder.UpdateColumns(&model.JPOrder{
					Status: model.JPOrderStatusInit,
				})
				if err != nil {
					utils.Errorf("UpdateColumns jporder(%d) error: %v", jporder.ID, err)
				}
				time.Sleep(time.Second * 1)
				continue
			}

			jadepoolOrderID, err := strconv.Atoi(result.ID)
			if err != nil {
				utils.Errorf("error: %v", err)
				err = jporder.UpdateColumns(&model.JPOrder{
					Status: model.JPOrderStatusInit,
				})
				if err != nil {
					utils.Errorf("UpdateColumns jporder(%d) error: %v", jporder.ID, err)
				}
				time.Sleep(time.Second * 1)
				continue
			}
			jporder.JadepoolOrderID = uint(jadepoolOrderID)
			jporder.From = result.From
			jporder.Confirmations = result.Confirmations
			jporder.Status = model.JPOrderStatusPending
			err = jporder.Save()
			if err != nil {
				utils.Errorf("error: %v", err)
				err = jporder.UpdateColumns(&model.JPOrder{
					Status: model.JPOrderStatusInit,
				})
				if err != nil {
					utils.Errorf("UpdateColumns jporder(%d) error: %v", jporder.ID, err)
				}
				time.Sleep(time.Second * 1)
				continue
			}

			utils.Infof("send jporder(%d) ok", jporder.ID)
			utils.Debugf("send jporder(%d) result: %v", jporder.ID, result)
		}

		time.Sleep(time.Second * 10)
	}
}
