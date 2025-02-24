package merchstore

type User struct {
	ID          int                `json:"-" db:"id"`
	Username    string             `json:"username" binding:"required"`
	Password    string             `json:"password" binding:"required"`
	Balance     int                `json:"coins" validate:"gte=0"`
	Inventory   []InventoryItem    `json:"inventory" validate:"dive"`
	CoinHistory TransactionHistory `json:"coinHistory" validate:"dive"`
}

type InventoryItem struct {
	ItemType string `json:"type" validate:"required"`
	Quantity int    `json:"quantity" validate:"gte=0"`
}

type Item struct {
	ID       int    `json:"-" validate:"-"`
	ItemType string `json:"type" validate:"required"`
	Price    int    `json:"-" validate:"gte=0"`
}

type TransactionHistory struct {
	Received []IncomingTransaction `json:"received" validate:"dive"`
	Sent     []OutgoingTransaction `json:"sent" validate:"dive"`
}

type IncomingTransaction struct {
	FromUser string `json:"fromUser" validate:"required"`
	Amount   int    `json:"amount" validate:"gt=0"`
}

type OutgoingTransaction struct {
	ToUser string `json:"toUser" validate:"required"`
	Amount int    `json:"amount" validate:"gt=0"`
}

type ErrorResponse struct {
	Message string `json:"errors" validate:"required"`
}

type AuthResponse struct {
	Token string `json:"token" validate:"required"`
}
