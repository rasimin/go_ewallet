package service

import (
	"context"
	"ewallet/gateaway/config"
	"ewallet/gateaway/model"
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

type Server struct {
	UserClient        pb.UserServiceClient
	TransactionClient pb.TransactionServiceClient
	UserID            uint32
}

func NewServer() *Server {
	// User connection
	connUser, err := grpc.NewClient(config.GetUserAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	// Transaction connection
	connTransaction, err := grpc.NewClient(config.GetTransactionAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return &Server{
		UserClient:        pb.NewUserServiceClient(connUser),
		TransactionClient: pb.NewTransactionServiceClient(connTransaction),
	}
}

func (s *Server) GetUserByID(c *gin.Context) {
	userIDParam := c.Param("userID")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := s.UserClient.GetUserByID(ctx, &pb.GetUserByIDRequest{UserId: uint32(userID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res.GetUser())
}

func (s *Server) GetWalletByUserID(c *gin.Context) {
	userIDParam := c.Param("userID")
	userID, err := strconv.ParseInt(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := s.TransactionClient.GetWalletByUserID(ctx, &pb.GetWalletByUserIDRequest{UserId: int32(userID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res.GetWallets())
}

func (s *Server) CreateUser(c *gin.Context) {
	var req pb.CreateUserRequest
	var reqW pb.CreateWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := s.UserClient.CreateUser(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	s.UserID = res.GetUser().UserId

	reqW.Wallet = &pb.Wallet{
		UserId:    s.UserID, // Menggunakan user ID yang baru dibuat
		Balance:   0,
		CreatedAt: timestamppb.New(time.Now()),
		UpdatedAt: timestamppb.New(time.Now()),
	}

	// Memanggil RPC CreateWallet
	_, err = s.TransactionClient.CreateWallet(ctx, &reqW)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res.GetUser())
}

func (s *Server) TransferWallet(c *gin.Context) {
	var req model.TransferWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	walletfrom, _ := s.TransactionClient.GetWalletByUserID(ctx, &pb.GetWalletByUserIDRequest{UserId: int32(req.UserIDFrom)})
	walletto, _ := s.TransactionClient.GetWalletByUserID(ctx, &pb.GetWalletByUserIDRequest{UserId: int32(req.UserIDTo)})

	res, err := s.TransactionClient.TransferWallet(ctx, &pb.TransferWalletRequest{
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

func (s *Server) TopUp(c *gin.Context) {
	var req model.TopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	wallet, err := s.TransactionClient.GetWalletByUserID(ctx, &pb.GetWalletByUserIDRequest{UserId: int32(req.UserID)})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	res, err := s.TransactionClient.TopUp(ctx, &pb.TopUpRequest{
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

func (s *Server) GetTransactionByUserID(c *gin.Context) {
	userIDParam := c.Param("userID")
	userID, err := strconv.ParseInt(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := s.TransactionClient.GetTransactionByUserID(ctx, &pb.GetTransactionByUserIDRequest{UserId: int32(userID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res.GetTransactions())
}

func (s *Server) GetUserAndBalanceWallet(c *gin.Context) {
	userIDParam := c.Param("userID")
	userID, err := strconv.ParseInt(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Fetch user details
	userRes, err := s.UserClient.GetUserByID(ctx, &pb.GetUserByIDRequest{UserId: uint32(userID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch wallet details
	walletRes, err := s.TransactionClient.GetWalletByUserID(ctx, &pb.GetWalletByUserIDRequest{UserId: int32(userID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Prepare combined response
	response := model.UserAndWalletResponse{
		User:   userRes.GetUser(),
		Wallet: walletRes.GetWallets(),
	}

	c.JSON(http.StatusOK, response)
}
