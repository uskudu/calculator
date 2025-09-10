package calculationService

// Calculation response structure
type Calculation struct {
	ID         string `gorm:"primaryKey" json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}

// CalculationRequest request structure
type CalculationRequest struct {
	Expression string `json:"expression"`
}
