package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/khris-xp/assessment-tax/common/dto"
	"github.com/khris-xp/assessment-tax/config"
	"github.com/khris-xp/assessment-tax/models"
	"github.com/khris-xp/assessment-tax/validate"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type PersonalDeductionController struct{}

func (p PersonalDeductionController) CreatePersonalDeduction(c echo.Context) error {
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
