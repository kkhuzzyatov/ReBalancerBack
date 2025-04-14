package entities

type CalcRequest struct {
	CurAllocation    []CurAllocation    `json:"cur_allocation"`
	TargetAllocation []TargetAllocation `json:"target_allocation"`
}

type CurAllocation struct {
	Ticker string `json:"ticker"`
	Number int    `json:"number"`
}

type TargetAllocation struct {
	Ticker  string  `json:"ticker"`
	Percent float64 `json:"percent"`
}
