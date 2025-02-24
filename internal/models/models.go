package models

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type AuthConfig struct {
	Salt      string
	SignedKey string
}

type Config struct {
	DB   DBConfig
	Auth AuthConfig
}

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type SendCoinRequest struct {
	Receiver string `json:"toUser" binding:"required"`
	Coins    int    `json:"amount" binding:"required"`
}

type UserInfoResponse struct {
	Balance     int                `json:"coins"`
	Inventory   []InventoryItem    `json:"inventory"`
	CoinHistory TransactionHistory `json:"coinHistory"`
}

type InventoryItem struct {
	ItemType string `json:"type"`
	Quantity int    `json:"quantity"`
}

type TransactionHistory struct {
	Received []IncomingTransaction `json:"received"`
	Sent     []OutgoingTransaction `json:"sent"`
}

type IncomingTransaction struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type OutgoingTransaction struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}
