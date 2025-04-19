package entities

type StockDB struct {
	ID          int
	UserID      int
	Ticker      string
	Number      int
	TargetShare float64
}
