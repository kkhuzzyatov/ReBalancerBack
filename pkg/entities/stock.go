package entities

type Stock struct {
	Ticker   string  `json:ticker`
	Lot      int     `json:"lot"`
	Price    float64 `json:"price"`
	AciValue float64 `json:"aci_value"`
}
