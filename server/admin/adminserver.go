package admin

import (
	"fmt"
	"strings"

	adminc "cybex-gateway/controller/admin"
	model "cybex-gateway/modeladmin"
	"cybex-gateway/server/middleware"
	"cybex-gateway/types"
	"cybex-gateway/utils/ecc"
	"cybex-gateway/utils/log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var tokens = map[string]bool{
	// "Bearer yangyu.DSADSDsadasd@dasd^YHN":  true,
	// "Bearer zhangyi.DDhhhCsadasd@dasd^YHN": true,
}

func authMiddleware(c *gin.Context) {
	tokenHeader := c.GetHeader("Authorization") //Grab the token from the header
	if tokens[tokenHeader] == false {           //Token is missing, returns with error code 403 Unauthorized
		c.AbortWithStatusJSON(400, gin.H{
			"message": "token不对",
		})
		return
	}
	log.Infoln(tokenHeader)
	// Call the next handler, which can be another middleware in the chain, or the final handler.
	c.Next()
}
func updateAssetsOne(c *gin.Context) {
	query := &model.Asset{}
	err := c.Bind(query)
	if query.GatewayAccount != "" || query.GatewayPass != "" {
		c.JSON(400, gin.H{
			"message": "不可更新网关账户，如非改不可请联系管理员",
		})
		return
	}
	if query.CYBName != "" {
		err = adminc.CheckCYB(query)
		if err != nil {
			log.Errorln("CheckCYB", err)
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}
	}
	if query.GatewayAccount != "" {
		err = adminc.CheckGateway(query)
		if err != nil {
			log.Errorln("CheckGateway", err)
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}
	}
	address, err := model.UpdateAsset(query)
	if err != nil {
		log.Errorln("GetAssets", err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, address)
}
func createAssetsOne(c *gin.Context) {
	query := &model.Asset{}
	err := c.Bind(query)
	err = (*query).ValidateCreate()
	if err != nil {
		log.Errorln("valid", err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = adminc.CheckAsset(query)
	if err != nil {
		log.Errorln("valid", err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	query.GatewayPass = query.GatewayPassword
	address, err := model.AssetsCreate(query)
	if err != nil {
		log.Errorln("GetAssets", err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, address)
}
func getAssets(c *gin.Context) {
	query := &model.Asset{}
	err := c.Bind(query)
	address, err := model.AssetsQuery(query)
	if err != nil {
		log.Errorln("GetAssets", err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, address)
}
func failOrder(c *gin.Context) {
	query := &model.JPOrder{}
	err := c.Bind(query)
	query2 := &model.JPOrder{}
	query2.ID = query.ID
	address, err := model.JPOrderFind(query2)
	if len(address) != 1 {
		c.JSON(400, gin.H{
			"message": "未找到id",
		})
		return
	}
	order := address[0]
	err = order.Failit()
	if err != nil {
		log.Errorln("failOrder", err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, address)
}
func getOrders(c *gin.Context) {
	query := &model.JPOrder{}
	err := c.Bind(query)
	if query.Limit == 0 {
		query.Limit = 20
	}
	address, total, err := model.OrderQuery(query)
	if err != nil {
		log.Errorln("getOrders", err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"data":  address,
		"total": total,
		"size":  query.Limit,
	})
}
func getAddress(c *gin.Context) {
	user := c.Param("user")
	asset := c.Param("asset")
	log.Infoln("GetAddress", user, asset)
	asset = strings.ToUpper(asset)
	address, err := adminc.GetAddress(user, asset)
	if err != nil {
		log.Errorln("user address", err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, address)
}
func verifyAddress(c *gin.Context) {
	address := c.Param("address")
	asset := c.Param("asset")
	asset = strings.ToUpper(asset)
	verifyResult, err := adminc.VerifyAddress(asset, address)
	if err != nil {
		log.Errorln("verifyResult", err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, verifyResult)
}
func newAddress(c *gin.Context) {
	user := c.Param("user")
	asset := c.Param("asset")
	log.Infoln("newAddress", user, asset)
	asset = strings.ToUpper(asset)
	address, err := adminc.NewAddress(user, asset)
	if err != nil {
		log.Errorln("user newAddress", err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, address)
}
func notDone(c *gin.Context) {
	interval := c.Param("interval")
	size := 20
	offset := 0
	res, err := adminc.RecordNotDone(interval, offset, size)
	if err != nil {
		log.Errorln("user address", err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, res)
}
func recordAssets(c *gin.Context) {
	user := c.Param("user")
	res, err := adminc.GetRecordAsset(user)
	if err != nil {
		log.Errorln("recordAssets", err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"records": res,
		"total":   len(res),
	})
}
func recordList(c *gin.Context) {
	user := c.Param("user")
	log.Infoln("GetRecord", user)
	query := &types.RecordsQuery{}
	err := c.Bind(query)
	if err != nil {
		log.Errorln("Bind", err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
	}
	query.FundType = strings.ToUpper(query.FundType)
	query.User = user
	if query.Size == "" {
		query.Size = "10"
	}
	if query.LastID == "" {
		query.LastID = "99999999"
	}
	if query.Offset == "" {
		query.Offset = "0"
	}
	log.Infoln("GetRecord", *query)
	res, total, err := adminc.GetRecord(query)
	var out []*types.Record
	for _, re1 := range res {
		confirms := fmt.Sprintf("%d", re1.Confirmations)
		record := &types.Record{
			Type:        re1.Type,
			ID:          re1.ID,
			UpdatedAt:   re1.UpdatedAt,
			CybexName:   re1.CybUser,
			OutAddr:     re1.OutAddr,
			Confirms:    confirms,
			Asset:       re1.Asset,
			OutHash:     re1.Hash,
			CybHash:     re1.CYBHash,
			TotalAmount: re1.TotalAmount.String(),
			Amount:      re1.Amount.String(),
			Fee:         re1.Fee.String(),
			Status:      re1.Status,
			CreatedAt:   re1.CreatedAt,
			Link:        re1.Link,
		}
		if re1.Type == "DEPOSIT" {
			record.GatewayAddr = re1.To
		}
		out = append(out, record)
	}
	if err != nil {
		log.Errorln("GetRecord", err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"records": out,
		"size":    query.Size,
		"total":   total,
	})
}

// StartServer ...
func StartServer() {
	adminc.InitNode()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"Content-Length"},
		// AllowCredentials: true,
	}))
	r.Use(middleware.RequestLogger())
	r.Use(middleware.GinBodyLogMiddleware)
	r.GET("/t", func(c *gin.Context) {
		ecc.TestECCSign()
		c.JSON(200, gin.H{})
	})
	usersigned := r.Group("/")
	if viper.GetBool("adminserver.auth") {
		usersigned.Use(authMiddleware)
	}
	usersigned.POST("/v1/assets/list", getAssets)
	usersigned.POST("/v1/assets/update", updateAssetsOne)
	usersigned.POST("/v1/assets/add", createAssetsOne)
	//
	usersigned.POST("/v1/orders/list", getOrders)
	usersigned.POST("/v1/orders/failed", failOrder)
	//
	usersigned.GET("/v1/record/undone/:interval", notDone)
	usersigned.GET("/v1/users/:user/assets/:asset/address", getAddress)
	usersigned.POST("/v1/users/:user/assets/:asset/address/new", newAddress)
	usersigned.GET("/v1/assets/:asset/address/:address/verify", verifyAddress)
	usersigned.GET("/v1/users/:user/records", recordList)
	usersigned.GET("/v1/users/:user/assets", recordAssets)
	port := viper.GetString("adminserver.port")
	log.Infoln("userserver start at", port)
	admintoken := viper.GetStringSlice("adminserver.tokens")
	for _, token := range admintoken {
		fulltoken := "Bearer " + token
		tokens[fulltoken] = true
	}
	r.Run(port) // listen and serve on 0.0.0.0:8080
}
