package libs

import (
	"github.com/khris-xp/assessment-tax/common/dto"
	"github.com/khris-xp/assessment-tax/types"
)

func CalculateTax(totalIncome float64, allowances []dto.AllowancesType) float64 {
	var totalAllowances float64
	for _, allowance := range allowances {
		if allowance.AllowancesType == "donation" && allowance.Amount > 100000 {
			totalAllowances = 100000
		} else {
			totalAllowances += allowance.Amount
		}
	}

	return totalIncome - totalAllowances
}

func CalculateTaxRate(income float64) types.TaxRangeType {
	if income >= 0 && income <= 150000 {
		return types.TaxRangeType{
			TaxIndex: 0,
			Income:   income,
		}
	}
	if income > 150000 && income <= 500000 {
		return types.TaxRangeType{
			TaxIndex: 1,
			Income:   (income - 150000) * 0.1,
		}
	}
	if income > 500000 && income <= 1000000 {
		return types.TaxRangeType{
			TaxIndex: 2,
			Income:   (income - 500000) * 0.15,
		}
	}
	if income > 1000000 && income <= 2000000 {
		return types.TaxRangeType{
			TaxIndex: 3,
			Income:   (income - 1000000) * 0.2,
		}
	}
	if income > 2000000 {
		return types.TaxRangeType{
			TaxIndex: 4,
			Income:   (income - 2000000) * 0.35,
		}
	}
	return types.TaxRangeType{
		TaxIndex: 0,
		Income:   0,
	}
}

func CalculateLevelTax(tax_rate int, index int, tax float64) float64 {
	if tax_rate == index {
		return tax
	} else {
		return 0
	}
}
