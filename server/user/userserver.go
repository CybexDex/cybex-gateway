package user

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/woyoutlz/bbb-gateway/server/middleware"
	"bitbucket.org/woyoutlz/bbb-gateway/types"
	"bitbucket.org/woyoutlz/bbb-gateway/utils/ecc"

	userc "bitbucket.org/woyoutlz/bbb-gateway/controller/user"
	"bitbucket.org/woyoutlz/bbb-gateway/utils/log"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func authMiddleware(c *gin.Context) {
	tokenHeader := c.GetHeader("Authorization") //Grab the token from the header
	if tokenHeader == "" {                      //Token is missing, returns with error code 403 Unauthorized
		c.AbortWithStatusJSON(400, gin.H{
			"message": "no token",
		})
		return
	}

	splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
	if len(splitted) != 2 {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "token format error",
		})
		return
	}

	tokenPart := splitted[1]
	tokenArr := strings.Split(tokenPart, ".")
	if len(tokenArr) != 3 {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "token format error 2",
		})
		return
	}
	signTime := tokenArr[0]
	signTimeInt, err := strconv.Atoi(signTime)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "signTimeInt error",
		})
		return
	}
	now := time.Now().Unix()
	if (int64(signTimeInt) < now-15*60) || (int64(signTimeInt) > now+15*60) {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "auth time error",
		})
		return
	}
	user := tokenArr[1]
	varsUser := c.Param("user")
	// check is tokenpart in db
	isok, _, err := userc.CheckUser(signTime, user, tokenArr[2])
	if err != nil {
		log.Errorln(err)
		c.AbortWithStatusJSON(400, gin.H{
			"message": "CheckUser err error",
		})
		return
	}
	if !isok {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "CheckUser fail",
		})
		return
	}
	if user == "" {
		// TODO uncomment this
		c.AbortWithStatusJSON(400, gin.H{
			"message": "no user error",
		})
		return
	}
	if varsUser != "" && varsUser != user {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "wrong user error",
		})
		return
	}
	// Call the next handler, which can be another middleware in the chain, or the final handler.
	c.Next()
}
func getAssets(c *gin.Context) {
	address, err := userc.GetAssets()
	if err != nil {
		log.Errorln("GetAssets", err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, address)
}
func getAddress(c *gin.Context) {
	user := c.Param("user")
	asset := c.Param("asset")
	log.Infoln("GetAddress", user, asset)
	asset = strings.ToUpper(asset)
	address, err := userc.GetAddress(user, asset)
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
	verifyResult, err := userc.VerifyAddress(asset, address)
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
	address, err := userc.NewAddress(user, asset)
	if err != nil {
		log.Errorln("user newAddress", err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, address)
}
func bbbAsset(c *gin.Context) {
	address, err := userc.GetBBBAssets()
	if err != nil {
		log.Errorln("user address", err)
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
	res, err := userc.RecordNotDone(interval, offset, size)
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
	res, err := userc.GetRecordAsset(user)
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
	log.Infoln("GetRecord", *query)
	res, total, err := userc.GetRecord(query)
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
	userc.InitNode()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(middleware.GinBodyLogMiddleware)
	r.GET("/t", func(c *gin.Context) {
		ecc.TestECCSign()
		c.JSON(200, gin.H{})
	})
	r.GET("/v1/assets", getAssets)
	r.GET("/v1/bbb", bbbAsset)
	r.GET("/v1/record/undone/:interval", notDone)
	usersigned := r.Group("/")
	if viper.GetBool("userserver.auth") {
		usersigned.Use(authMiddleware)
	}
	usersigned.GET("/v1/users/:user/assets/:asset/address", getAddress)
	usersigned.POST("/v1/users/:user/assets/:asset/address/new", newAddress)
	usersigned.GET("/v1/assets/:asset/address/:address/verify", verifyAddress)
	usersigned.GET("/v1/users/:user/records", recordList)
	usersigned.GET("/v1/users/:user/assets", recordAssets)

	port := viper.GetString("userserver.port")
	log.Infoln("userserver start at", port)
	r.Run(port) // listen and serve on 0.0.0.0:8080
}
