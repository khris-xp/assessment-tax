package controllers

import (
	"net/http"

	"github.com/khris-xp/assessment-tax/common/dto"
	"github.com/khris-xp/assessment-tax/constants"
	"github.com/khris-xp/assessment-tax/libs"
	"github.com/khris-xp/assessment-tax/response"
	"github.com/khris-xp/assessment-tax/validate"
	"github.com/labstack/echo/v4"
)

type TaxController struct{}

func (t TaxController) CalculateTax(c echo.Context) error {
	var taxRequest dto.TaxRequest

	if err := c.Bind(&taxRequest); err != nil {
		return c.JSON(http.StatusBadRequest, "Expected bad request error, got no error or different status code")
	}

	if err := validate.ValidateTaxRequest(taxRequest); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	taxLevels := []struct {
		Range string
		Min   int
		Max   int
	}{
		{"0-150,000", 0, 150000},
		{"150,001-500,000", 150001, 500000},
		{"500,001-1,000,000", 500001, 1000000},
		{"1,000,001-2,000,000", 1000001, 2000000},
		{"2,000,001 ขึ้นไป", 2000001, -1},
	}

	tax_calculate := libs.CalculateTax(taxRequest.TotalIncome, taxRequest.Allowances)
	tax_total_with_deduction := tax_calculate - constants.TaxDeductionInit().Deduction

	tax_rate := libs.CalculateTaxRate(tax_total_with_deduction)
	tax := float64(tax_rate.Income) - taxRequest.Wht

	var taxLevelResponses []response.TaxLevelResponse
	for index, level := range taxLevels {
		taxLevelResponses = append(taxLevelResponses, response.TaxLevelResponse{
			Level: level.Range,
			Tax:   libs.CalculateLevelTax(tax_rate.TaxIndex, index, tax),
		})
	}

	response := map[string]interface{}{
		"tax":      tax,
		"taxLevel": taxLevelResponses,
	}

	return c.JSON(http.StatusOK, response)

}
