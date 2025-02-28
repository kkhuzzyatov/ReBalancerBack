package entities

type User struct {
	ID               int    `json:"id"`
	Email            string `json:"email"`
	PasswordHash     string `json:"password"`
	CurAllocation    string `json:"curAllocation"`
	TargetAllocation string `json:"targetAllocation"`
	TaxRate          int    `json:"taxRate"`
}