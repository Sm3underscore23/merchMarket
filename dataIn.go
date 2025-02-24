package merchstore

type SendCoinRequest struct {
	Receiver string `json:"toUser" validate:"required"`
	Coins    int    `json:"amount" validate:"gt=0"`
}

type AuthRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
