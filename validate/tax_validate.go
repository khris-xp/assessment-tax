package validate

import (
	"errors"

	"github.com/khris-xp/assessment-tax/common/dto"
)

func ValidateTaxRequest(req dto.TaxRequest) error {
	if req.TotalIncome < 0 {
		return errors.New("Total income must be greater than or equal to 0")
	}

	if len(req.Allowances) == 0 {
		return errors.New("Allowances must not be empty")
	}

	if req.Wht > req.TotalIncome {
		return errors.New("Withholding tax must be less than total income")
	}

	if req.TotalIncome <= 150000 {
		return nil
	}

	if (req.Allowances[0].AllowancesType == "donation") && (req.Allowances[0].Amount < 0) {
		return errors.New("Donation amount must be greater than or equal to 0")
	}

	return nil
}

func ValidatePersonalDeductionAmount(amount float64) error {
	if amount < 60000 {
		return errors.New("Personal deduction must be at least 60,000")
	}
	if amount > 100000 {
		return errors.New("Personal deduction must be at most 100,000")
	}
	return nil
}

func ValidateKReceiptAmount(amount float64) error {
	if amount < 60000 {
		return errors.New("K-Receipt must be at least 50,000")
	}
	if amount > 100000 {
		return errors.New("K-Receipt must be at most 100,000")
	}
	return nil
}
