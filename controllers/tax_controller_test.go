package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
)

func TestCalculateTax(t *testing.T) {
	tests := []struct {
		name        string
		totalIncome float64
		wht         float64
		allowances  []struct {
			allowanceType string
			amount        float64
		}
		statusCode int
	}{
		{
			"StatusOK with allowances",
			70000.0,
			2000.0,
			[]struct {
				allowanceType string
				amount        float64
			}{
				{
					"donation",
					0.0,
				},
			},
			http.StatusOK,
		},
		{"StatusOK when Income 70,000", 70000.0, 0.0, []struct {
			allowanceType string
			amount        float64
		}{
			{
				"donation",
				0.0,
			},
		}, http.StatusOK},
		{"StatusOK when Income 150,000", 150000.0, 0.0, []struct {
			allowanceType string
			amount        float64
		}{
			{
				"donation",
				0.0,
			},
		}, http.StatusOK},
		{"StatusOK when Income 300,000", 300000.0, 0.0, []struct {
			allowanceType string
			amount        float64
		}{
			{
				"donation",
				0.0,
			},
		}, http.StatusOK},

		{"StatusOK when Income 70,000, WHT 10,000", 70000.0, 10000.0, []struct {
			allowanceType string
			amount        float64
		}{
			{
				"donation",
				0.0,
			},
		}, http.StatusOK},
		{"StatusOK when Income 150,000, WHT 10,000", 150000.0, 10000.0, []struct {
			allowanceType string
			amount        float64
		}{
			{
				"donation",
				0.0,
			},
		}, http.StatusOK},

		{"StatusOK when Income MaxFloat64, WHT MaxFloat64", 1.7976931348623157e+308, 1.7976931348623157e+308, []struct {
			allowanceType string
			amount        float64
		}{
			{
				"donation",
				0.0,
			},
		}, http.StatusOK},
		{"StatusOK when Income 0.5, WHT 0.1", 0.5, 0.1, []struct {
			allowanceType string
			amount        float64
		}{
			{
				"donation",
				0.0,
			},
		}, http.StatusOK},

		{"StatusOK when Income 70,000, WHT 10,000, Donation 100,001", 70000.0, 10000.0, []struct {
			allowanceType string
			amount        float64
		}{
			{
				"donation",
				100001.0,
			},
		}, http.StatusOK},
		{"StatusOK  when Income 70,000, WHT 10,000, Donation 100,000", 70000.0, 10000.0, []struct {
			allowanceType string
			amount        float64
		}{
			{
				"donation",
				100000.0,
			},
		}, http.StatusOK},

		{"StatusBadRequest when Income MinFloat64, WHT MinFloat64", -1.7976931348623157e+308, -1.7976931348623157e+308, []struct {
			allowanceType string
			amount        float64
		}{
			{
				"donation",
				0.0,
			},
		}, http.StatusBadRequest},

		{"StatusBadRequest when Income -10", -10.0, 0.0, []struct {
			allowanceType string
			amount        float64
		}{
			{
				"donation",
				0.0,
			},
		}, http.StatusBadRequest},
		{"StatusBadRequest when WHT > Income", 10000.0, 20000.0, []struct {
			allowanceType string
			amount        float64
		}{
			{
				"donation",
				0.0,
			},
		}, http.StatusBadRequest},
		{"StatusBadRequest when Allowances empty", 10000.0, 0.0, []struct {
			allowanceType string
			amount        float64
		}{}, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			var allowancesJSON string
			if len(tt.allowances) > 0 {
				var allowanceStrings []string
				for _, allowance := range tt.allowances {
					allowanceStrings = append(allowanceStrings, fmt.Sprintf(`{"allowanceType": "%s", "amount": %v}`, allowance.allowanceType, allowance.amount))
				}
				allowancesJSON = fmt.Sprintf(`"allowances": [%s], `, strings.Join(allowanceStrings, ","))
			}

			reqBody := fmt.Sprintf(`{%s"totalIncome": %v, "wht": %v}`, allowancesJSON, tt.totalIncome, tt.wht)
			req := httptest.NewRequest(http.MethodPost, "/tax/calculations", strings.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			tCl := TaxController{}
			err := tCl.CalculateTax(c)

			if (err != nil && tt.statusCode != http.StatusBadRequest) || rec.Code != tt.statusCode {
				fmt.Printf(color.RedString("%s failed. Expected status code: %v, got: %v\n"), tt.name, tt.statusCode, rec.Code)
				return
			}

			fmt.Printf(color.GreenString("%s passed. Income: %v, WHT: %v\n"), tt.name, tt.totalIncome, tt.wht)
		})
	}
}