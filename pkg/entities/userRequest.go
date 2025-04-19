package entities

type UserRequest struct {
	Email            string             `json:"email"`
	Password         string             `json:"password"`
	CurAllocation    []CurAllocation    `json:"cur_allocation"`
	TargetAllocation []TargetAllocation `json:"target_allocation"`
}
