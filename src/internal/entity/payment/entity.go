package payment

type Bill struct {
	TgID  int64   `json:"tg_id"`
	Price float64 `json:"price"`
	SubID int     `json:"sub_id"`
}

func NewBill(
	tgID int64,
	price float64,
	subID int,
) Bill {
	return Bill{
		TgID:  tgID,
		Price: price,
		SubID: subID,
	}
}
