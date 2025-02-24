package bot

//Возвращает map[тикер]процент
func GetTargetAllocation() (map[string]float64) {
	targetAllocation := make(map[string]float64)
	targetAllocation["TMOS"] = 60
	targetAllocation["TBRU"] = 36.5
	targetAllocation["RUB"] = 3.5
	return targetAllocation
}

//Возвращает map[тикер]штуки
func GetCurAllocation() (map[string]int) {
	curAllocation := make(map[string]int)
	curAllocation["TMOS"] = 12000
	curAllocation["TBRU"] = 8000
	curAllocation["RUB"] = 4500
	return curAllocation
}

//Возвращает сумму пополнения в рублях
func GetReplenishmentAmount() (int) {
	return 5000
}

//Возвращает налоговую ставку
func GetTaxRate() (float64) {
	return 13
}