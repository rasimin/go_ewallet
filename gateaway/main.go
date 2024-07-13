package main

import (
	"context"
	pb "ewallet/gateaway/proto"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	userAddress        = "localhost:50052"
	transactionAddress = "localhost:50051"
)

func main() {
	// user
	connuser, err := grpc.NewClient(userAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer connuser.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// transaction
	conntransaction, err := grpc.NewClient(transactionAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conntransaction.Close()
	ctx, canceltrans := context.WithTimeout(context.Background(), time.Second)
	defer canceltrans()

	c := pb.NewUserServiceClient(connuser)
	d := pb.NewTransactionServiceClient(conntransaction)
	userID := uint32(1) // replace with the user ID you want to fetch

	r, err := c.GetUserByID(ctx, &pb.GetUserByIDRequest{UserId: userID})
	if err != nil {
		log.Fatalf("could not get user: %v", err)
	}
	fmt.Printf("User: %v\n", r.GetUser())

	userIDx := int32(1)

	t, err := d.GetWalletByUserID(ctx, &pb.GetWalletByUserIDRequest{UserId: userIDx})
	if err != nil {
		log.Fatalf("could not get user: %v", err)
	}
	fmt.Printf("User: %v\n", t.GetWallets())

}
