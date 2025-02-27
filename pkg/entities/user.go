package entities

type User struct {
	userID           int64
	curAllocation    map[string]int
	targetAllocation map[string]float64
	taxRate          float64
}

func NewUser(userID int64, taxRate float64) *User {
	return &User{
		userID:           userID,
		curAllocation:    make(map[string]int),
		targetAllocation: make(map[string]float64),
		taxRate:          taxRate,
	}
}

func (u *User) GetUserID() int64 {
	return u.userID
}

func (u *User) GetCurAllocation() map[string]int {
	copy := make(map[string]int)
	for k, v := range u.curAllocation {
		copy[k] = v
	}
	return copy
}

func (u *User) GetTargetAllocation() map[string]float64 {
	copy := make(map[string]float64)
	for k, v := range u.targetAllocation {
		copy[k] = v
	}
	return copy
}

func (u *User) GetTaxRate() float64 {
	return u.taxRate
}

func (u *User) SetCurAllocation(allocation map[string]int) {
	u.curAllocation = make(map[string]int)
	for k, v := range allocation {
		u.curAllocation[k] = v
	}
}

func (u *User) SetTargetAllocation(allocation map[string]float64) {
	u.targetAllocation = make(map[string]float64)
	for k, v := range allocation {
		u.targetAllocation[k] = v
	}
}

func (u *User) SetTaxRate(rate float64) {
	u.taxRate = rate
}