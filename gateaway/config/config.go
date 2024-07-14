// config/config.go
package config

const (
	UserAddress        = "localhost:50052"
	TransactionAddress = "localhost:50051"
	HTTPPort           = ":8080"
)

func GetUserAddress() string {
	return UserAddress
}

func GetTransactionAddress() string {
	return TransactionAddress
}

func GetHTTPPort() string {
	return HTTPPort
}
