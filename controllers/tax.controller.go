package controllers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gocarina/gocsv"
	"github.com/khris-xp/assessment-tax/common/dto"
	"github.com/khris-xp/assessment-tax/config"
	"github.com/khris-xp/assessment-tax/constants"
	"github.com/khris-xp/assessment-tax/libs"
	"github.com/khris-xp/assessment-tax/models"
	"github.com/khris-xp/assessment-tax/response"
	"github.com/khris-xp/assessment-tax/validate"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type TaxController struct{}

func (t TaxController) CalculateTax(c echo.Context) error {
	var taxRequest dto.TaxRequest
	db := config.DB()
	var personal_deduction models.PersonalDeduction

	if err := db.First(&personal_deduction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusBadRequest, "Personal deduction is not set")
		} else {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to get personal deduction: %v", err))
		}
	}

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
	tax_total_with_deduction := tax_calculate - personal_deduction.PersonalDeduction

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

func (t TaxController) CreatePersonalDeduction(c echo.Context) error {
	var personal_deduction dto.PersonalDeductionRequest
	if err := c.Bind(&personal_deduction); err != nil {
		return c.JSON(http.StatusBadRequest, "Expected bad request error, got no error or different status code")
	}

	if err := validate.ValidatePersonalDeductionAmount(personal_deduction.Amount); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	db := config.DB()

	var existingPersonalDeduction models.PersonalDeduction
	if err := db.First(&existingPersonalDeduction).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to check existing personal deduction: %v", err))
	}

	if existingPersonalDeduction.ID != 0 {
		existingPersonalDeduction.PersonalDeduction = personal_deduction.Amount
		if err := db.Save(&existingPersonalDeduction).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to update personal deduction: %v", err))
		}
	} else {
		newPersonalDeduction := models.PersonalDeduction{PersonalDeduction: personal_deduction.Amount}
		if err := db.Create(&newPersonalDeduction).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to create personal deduction: %v", err))
		}
	}

	response := map[string]interface{}{
		"personalDeduction": personal_deduction.Amount,
	}
	return c.JSON(http.StatusOK, response)
}

func (c *TaxController) TaxCalculateFormCsv(ctx echo.Context) error {
	taxFile, err := ctx.FormFile("taxFile")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Expected bad request error, got no error or different status code")
	}
	fileCsv, err := taxFile.Open()

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to open file: %v", err))
	}

	defer fileCsv.Close()

	taxData, _ := io.ReadAll(fileCsv)

	var taxCsv []dto.TaxCSV
	if err = gocsv.UnmarshalBytes(taxData, &taxCsv); err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to unmarshal csv: %v", err))
	}

	var taxes []response.TaxResult

	for _, tax := range taxCsv {
		taxRequest := dto.TaxRequest{
			TotalIncome: tax.TotalIncome,
			Wht:         tax.Wht,
			Allowances: []dto.AllowancesType{
				{
					AllowancesType: "donation",
					Amount:         tax.Donation,
				},
			},
		}

		if err := validate.ValidateTaxRequest(taxRequest); err != nil {
			return ctx.JSON(http.StatusBadRequest, err.Error())
		}

		tax_calculate := libs.CalculateTax(taxRequest.TotalIncome, taxRequest.Allowances)
		tax_total_with_deduction := tax_calculate - constants.TaxDeductionInit().Deduction
		tax_rate := libs.CalculateTaxRate(tax_total_with_deduction)
		tax := float64(tax_rate.Income) - taxRequest.Wht

		if tax < 0 {
			tax = 0
		}

		taxResult := response.TaxResult{
			TotalIncome: taxRequest.TotalIncome,
			Tax:         tax,
		}

		taxes = append(taxes, taxResult)
	}
	response := map[string]interface{}{
		"taxes": taxes,
	}

	return ctx.JSON(http.StatusOK, response)
}
