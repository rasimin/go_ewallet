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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TransferWalletRequest struct {
	UserIDFrom int32   `json:"user_idfrom"`
	UserIDTo   int32   `json:"user_idto"`
	Amount     float32 `json:"amount"`
}

type UserAndWalletResponse struct {
	User   *pb.User   `json:"user"`
	Wallet *pb.Wallet `json:"wallet"`
}

const (
	userAddress        = "localhost:50052"
	transactionAddress = "localhost:50051"
)

type server struct {
	userClient        pb.UserServiceClient
	transactionClient pb.TransactionServiceClient
	userid            uint32
}

type TopUpRequest struct {
	UserID int32   `json:"user_id"`
	Amount float32 `json:"amount"`
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

func (s *server) createUser(c *gin.Context) {
	var req pb.CreateUserRequest
	var reqW pb.CreateWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := s.userClient.CreateUser(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	s.userid = res.GetUser().UserId

	reqW.Wallet = &pb.Wallet{
		UserId:    s.userid, // Menggunakan user ID yang baru dibuat
		Balance:   0,
		CreatedAt: timestamppb.New(time.Now()),
		UpdatedAt: timestamppb.New(time.Now()),
	}

	// Memanggil RPC CreateWallet
	_, err = s.transactionClient.CreateWallet(ctx, &reqW)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res.GetUser())
}

func (s *server) transferWallet(c *gin.Context) {
	var req TransferWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	walletfrom, _ := s.transactionClient.GetWalletByUserID(ctx, &pb.GetWalletByUserIDRequest{UserId: int32(req.UserIDFrom)})
	walletto, _ := s.transactionClient.GetWalletByUserID(ctx, &pb.GetWalletByUserIDRequest{UserId: int32(req.UserIDTo)})

	res, err := s.transactionClient.TransferWallet(ctx, &pb.TransferWalletRequest{
		FromWalletId: walletfrom.Wallets.Id,
		ToWalletId:   walletto.Wallets.Id,
		Amount:       req.Amount,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": res.Message})
}
func (s *server) topUp(c *gin.Context) {
	var req TopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	wallet, err := s.transactionClient.GetWalletByUserID(ctx, &pb.GetWalletByUserIDRequest{UserId: int32(req.UserID)})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	res, err := s.transactionClient.TopUp(ctx, &pb.TopUpRequest{
		WalletId: wallet.Wallets.Id,
		Amount:   req.Amount,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Top-up successful",
		"transaction": res.Transaction,
	})
}

func (s *server) getTransactionByUserID(c *gin.Context) {
	userIDParam := c.Param("userID")
	userID, err := strconv.ParseInt(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := s.transactionClient.GetTransactionByUserID(ctx, &pb.GetTransactionByUserIDRequest{UserId: int32(userID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res.GetTransactions())
}

func (s *server) getUserAndBalanceWallet(c *gin.Context) {
	userIDParam := c.Param("userID")
	userID, err := strconv.ParseInt(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Fetch user details
	userRes, err := s.userClient.GetUserByID(ctx, &pb.GetUserByIDRequest{UserId: uint32(userID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch wallet details
	walletRes, err := s.transactionClient.GetWalletByUserID(ctx, &pb.GetWalletByUserIDRequest{UserId: int32(userID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Prepare combined response
	response := UserAndWalletResponse{
		User:   userRes.GetUser(),
		Wallet: walletRes.GetWallets(),
	}

	c.JSON(http.StatusOK, response)
}

func main() {
	srv := newServer()

	r := gin.Default()
	r.GET("/getUserByID/:userID", srv.getUserByID)
	r.GET("/getWalletByUserID/:userID", srv.getWalletByUserID)
	r.GET("/getTransactionByUserID/:userID", srv.getTransactionByUserID)
	r.GET("/getUserAndBalanceWallet/:userID", srv.getUserAndBalanceWallet)

	r.POST("/createUser", srv.createUser)
	r.POST("/transferWallet", srv.transferWallet)
	r.POST("/topUp", srv.topUp)

	log.Println("Starting HTTP server on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
