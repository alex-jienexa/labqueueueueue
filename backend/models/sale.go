package models

// Sale предоставляет объект ситуации продажи места в очереди
// за определённую сумму (в рублях). Хранит продавца, покупателя и
// факт подтверждения сделки.
type Sale struct {
	ID        int  `json:"id"`
	QueueID   int  `json:"queue_id"`
	SellerID  int  `json:"seller_id"`
	BuyerID   int  `json:"buyer_id"`
	Price     int  `json:"price"`
	Confirmed bool `json:"confirmed"`
}
