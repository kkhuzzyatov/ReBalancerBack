package entities

type CalcResponse struct {
	StockData      []Stock    `json:"stock_data"`
	TotalValue     float64    `json:"total_value"`
	Orders         []Order    `json:"orders"`
	FinalStructure []Position `json:"final_structure"`
	Err            string     `json:"err"`
}

type Order struct {
	Ticker    string `json:"ticker"`
	NumberLot int    `json:"number_lot"`
}

type Position struct {
	Ticker           string  `json:"ticker"`
	Number           float64 `json:"number"`
	PercentOfCapital float64 `json:"percentOfCapital"`
}
