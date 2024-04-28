package dto

type TaxRequest struct {
	TotalIncome float64          `json:"totalIncome"`
	Wht         float64          `json:"wht"`
	Allowances  []AllowancesType `json:"allowances"`
}

type AllowancesType struct {
	AllowancesType string  `json:"allowanceType"`
	Amount         float64 `json:"amount"`
}

type TaxCSV struct {
	TotalIncome float64 `csv:"totalIncome"`
	Donation    float64 `csv:"donation"`
	Wht         float64 `csv:"wht"`
}
