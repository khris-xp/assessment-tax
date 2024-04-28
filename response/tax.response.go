package response

type TaxResponse struct {
	Tax      float64            `json:"tax"`
	TaxLevel []TaxLevelResponse `json:"taxLevel"`
}

type TaxLevelResponse struct {
	Level string     `json:"level"`
	Tax   float64 `json:"tax"`
}

type TaxResult struct {
	TotalIncome float64 `json:"totalIncome"`
	Tax         float64 `json:"tax"`
}