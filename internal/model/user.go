package model

type User struct {
	ID    int `json:"id"`
	Coins int `json:"coins"`
}

type Wallet struct {
	Coins              int                  `json:"coins"`
	TransactionHistory []WalletHistoryEntry `json:"transaction_history"`
}

type WalletHistoryEntry struct {
	TransactionType string `json:"transaction_type"` 
	CounterpartyID  int    `json:"counterparty_id"`
	Amount          int    `json:"amount"`
	CreatedAt       string `json:"created_at"` 
}
