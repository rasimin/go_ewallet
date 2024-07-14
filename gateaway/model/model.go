// model/model.go
package model

import (
	pb "ewallet/gateaway/proto"
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

type TopUpRequest struct {
	UserID int32   `json:"user_id"`
	Amount float32 `json:"amount"`
}
