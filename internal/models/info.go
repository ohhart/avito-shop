package models

// История переводов монет
type CoinHistory struct {
	Received []TransactionHistory `json:"received"`
	Sent     []TransactionHistory `json:"sent"`
}

// Описание одной транзакции в истории
type TransactionHistory struct {
	FromUser string `json:"fromUser,omitempty"`
	ToUser   string `json:"toUser,omitempty"`
	Amount   int    `json:"amount"`
}

// Структура для /api/info
type InventoryResponse struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type InfoResponse struct {
	Coins       int                 `json:"coins"`
	Inventory   []InventoryResponse `json:"inventory"`
	CoinHistory CoinHistory         `json:"coinHistory"`
}
