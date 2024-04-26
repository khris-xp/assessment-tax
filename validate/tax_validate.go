package validate

import (
	"errors"
	"fmt"

	"github.com/khris-xp/assessment-tax/common/dto"
)

func ValidateTaxRequest(req dto.TaxRequest) error {
	if req.TotalIncome < 0 || req.Wht < 0 {
		return errors.New("Total income and wht must be greater than or equal to 0")
	}

	if len(req.Allowances) == 0 {
		return errors.New("Allowances must not be empty")
	}

	if req.Wht > req.TotalIncome {
		return errors.New("Withholding tax must be less than total income")
	}

	if (req.Allowances[0].AllowancesType == "donation") && (req.Allowances[0].Amount < 0) {
		fmt.Println("Error : ", "Donation amount must be greater than or equal to 0")
		return errors.New("Donation amount must be greater than or equal to 0")
	}

	if req.TotalIncome <= 150000 {
		return nil
	}

	return nil
}
