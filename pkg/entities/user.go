package entities

type User struct {
	CurAllocation    string  `json:"cur_allocation"`
	TargetAllocation string  `json:"target_allocation"`
	Cash             float64 `json:"cash"`
}