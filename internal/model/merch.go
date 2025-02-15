package model

type Merch struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type Purchase struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	ItemName    string `json:"item_name"`
	Price       int    `json:"price"`
	PurchasedAt string `json:"purchased_at"`
}
