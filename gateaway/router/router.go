// router/router.go
package router

import (
	"ewallet/gateaway/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter(srv *service.Server) *gin.Engine {
	r := gin.Default()
	r.GET("/getUserByID/:userID", srv.GetUserByID)
	r.GET("/getWalletByUserID/:userID", srv.GetWalletByUserID)
	r.GET("/getTransactionByUserID/:userID", srv.GetTransactionByUserID)
	r.GET("/getUserAndBalanceWallet/:userID", srv.GetUserAndBalanceWallet)

	r.POST("/createUser", srv.CreateUser)
	r.POST("/transferWallet", srv.TransferWallet)
	r.POST("/topUp", srv.TopUp)

	return r
}
