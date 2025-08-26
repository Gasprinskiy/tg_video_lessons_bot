package subscritions

type Subscrition struct {
	ID          int     `db:"sub_id"`
	TermInMonth int     `db:"term_in_month"`
	Price       float64 `db:"price"`
}
