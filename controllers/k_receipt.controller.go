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

type KReceiptController struct{}

func (k KReceiptController) CreateKReceipt(c echo.Context) error {
	var k_receipt dto.KReceiptRequest
	if err := c.Bind(&k_receipt); err != nil {
		return c.JSON(http.StatusBadRequest, "Expected bad request error, got no error or different status code")
	}

	if err := validate.ValidateKReceiptAmount(k_receipt.Amount); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	db := config.DB()

	var existingKReceipt models.KReceipt
	if err := db.First(&existingKReceipt).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to check existing K receipt: %v", err))
	}

	if existingKReceipt.ID != 0 {
		existingKReceipt.KReceipt = k_receipt.Amount
		if err := db.Save(&existingKReceipt).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to update K receipt: %v", err))
		}
	} else {
		newKReceipt := models.KReceipt{KReceipt: k_receipt.Amount}
		if err := db.Create(&newKReceipt).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to create K receipt: %v", err))
		}
	}

	response := map[string]interface{}{
		"kReceipt": k_receipt.Amount,
	}
	return c.JSON(http.StatusOK, response)
}