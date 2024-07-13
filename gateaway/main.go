package main

import (
	"context"
	pb "ewallet/gateaway/proto"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	userAddress        = "localhost:50052"
	transactionAddress = "localhost:50051"
)

type userRequest struct {
	UserID uint32 `json:"user_id"`
}
type server struct {
	userClient        pb.UserServiceClient
	transactionClient pb.TransactionServiceClient
}
type walletRequest struct {
	UserID int32 `json:"user_id"`
}

func newServer() *server {
	// User connection
	connUser, err := grpc.NewClient(userAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	// Transaction connection
	connTransaction, err := grpc.NewClient(transactionAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return &server{
		userClient:        pb.NewUserServiceClient(connUser),
		transactionClient: pb.NewTransactionServiceClient(connTransaction),
	}
}
func (s *server) getUserByID(c *gin.Context) {
	userIDParam := c.Param("userID")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := s.userClient.GetUserByID(ctx, &pb.GetUserByIDRequest{UserId: uint32(userID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res.GetUser())
}

func (s *server) getWalletByUserID(c *gin.Context) {
	userIDParam := c.Param("userID")
	userID, err := strconv.ParseInt(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := s.transactionClient.GetWalletByUserID(ctx, &pb.GetWalletByUserIDRequest{UserId: int32(userID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res.GetWallets())
}

func main() {
	srv := newServer()

	r := gin.Default()
	r.GET("/getUserByID/:userID", srv.getUserByID)
	r.GET("/getWalletByUserID/:userID", srv.getWalletByUserID)

	log.Println("Starting HTTP server on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
