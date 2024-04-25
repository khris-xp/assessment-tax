package controllers

import (
	"net/http"

	"github.com/khris-xp/assessment-tax/common/dto"
	"github.com/khris-xp/assessment-tax/constants"
	"github.com/khris-xp/assessment-tax/libs"
	"github.com/khris-xp/assessment-tax/response"
	"github.com/labstack/echo/v4"
)

type TaxController struct{}

func (t TaxController) CalculateTax(c echo.Context) error {

	var taxRequest dto.TaxRequest
	if err := c.Bind(&taxRequest); err != nil {
		return c.JSON(http.StatusBadRequest, "Expected bad request error, got no error or different status code")
	}

	if taxRequest.TotalIncome < 0 {
		return c.JSON(http.StatusBadRequest, "Total income must be greater than or equal to 0")
	}

	if len(taxRequest.Allowances) == 0 {
		return c.JSON(http.StatusBadRequest, "Allowances must not be empty")
	}

	if taxRequest.Wht > taxRequest.TotalIncome {
		return c.JSON(http.StatusBadRequest, "Withholding tax must be less than total income")
	}

	if taxRequest.TotalIncome <= 150000 {
		return c.JSON(http.StatusOK, response.TaxResponse{Tax: 0})
	}

	tax_calculate := libs.CalculateTax(taxRequest.TotalIncome, taxRequest.Allowances)
	tax_total_with_deduction := tax_calculate - constants.TaxDeductionInit().Deduction
	tax_rate := libs.CalculateTaxRate(tax_total_with_deduction)
	tax := tax_rate - taxRequest.Wht
	return c.JSON(http.StatusOK, response.TaxResponse{Tax: tax})
}
