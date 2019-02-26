package main

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"

	"git.coding.net/bobxuyang/cy-gateway-BN/models"
	"git.coding.net/bobxuyang/cy-gateway-BN/repository/jporder"

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
		w.Write([]byte(`{"status": true}`))
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
			data, err := controllers.DoSendOrder(jporder.ID)
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
				time.Sleep(time.Second * 1)
				continue
			}
			result := controllers.OrderNotiResult{}
			err = json.Unmarshal(resultBytes, &result)
			if err != nil {
				utils.Errorf("error: %v", err)
				time.Sleep(time.Second * 1)
				continue
			}

			jadepoolOrderID, err := strconv.Atoi(result.ID)
			if err != nil {
				utils.Errorf("error: %v", err)
				time.Sleep(time.Second * 1)
				continue
			}
			jporder.JadepoolOrderID = uint(jadepoolOrderID)
			jporder.From = result.From
			jporder.Confirmations = result.Confirmations
			err = jporder.Save()
			if err != nil {
				utils.Errorf("error: %v", err)
				time.Sleep(time.Second * 1)
				continue
			}

			utils.Infof("send jporder(%d) ok", jporder.ID)
			utils.Debugf("send jporder(%d) result: %v", jporder.ID, result)
		}

		time.Sleep(time.Second * 10)
	}
}
