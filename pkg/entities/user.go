package entities

type User struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	CurAllocation    string `json:"curAllocation"`
	TargetAllocation string `json:"targetAllocation"`
	TaxRate          int    `json:"taxRate"`
}