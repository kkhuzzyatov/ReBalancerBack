package entities

type User struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	CurAllocation    string `json:"cur_allocation"`
	TargetAllocation string `json:"target_allocation"`
}