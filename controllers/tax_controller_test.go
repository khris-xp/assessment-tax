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
		statusCode  int
	}{
		{"StatusBadRequest when Income -10", -10.0, 0.0, http.StatusBadRequest},
		{"StatusOK when Income 0", 0.0, 0.0, http.StatusOK},
		{"StatusOK when Income 70,000", 70000.0, 0.0, http.StatusOK},
		{"StatusOK when Income 150,000", 150000.0, 0.0, http.StatusOK},
		{"StatusOK when Income 300,000", 300000.0, 0.0, http.StatusOK},

		{"StatusOK when Income 70,000, WHT 10,000", 70000.0, 10000.0, http.StatusOK},
		{"StatusOK when Income 150,000, WHT 10,000", 150000.0, 10000.0, http.StatusOK},

		{"StatusOK when Income MaxFloat64, WHT MaxFloat64", 1.7976931348623157e+308, 1.7976931348623157e+308, http.StatusOK},
		{"StatusOK when Income MinFloat64, WHT MinFloat64", -1.7976931348623157e+308, -1.7976931348623157e+308, http.StatusBadRequest},

		{"StatusOK when Income 0.5, WHT 0.1", 0.5, 0.1, http.StatusOK},
		{"StatusBadRequest when WHT > Income", 10000.0, 20000.0, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			reqBody := fmt.Sprintf(`{"totalIncome": %v, "wht": %v}`, tt.totalIncome, tt.wht)
			req := httptest.NewRequest(http.MethodPost, "/tax/calculations", strings.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			tCl := TaxController{}
			err := tCl.CalculateTax(c)

			if (err != nil && tt.statusCode != http.StatusBadRequest) || rec.Code != tt.statusCode {
				t.Errorf("%s failed. Income: %v, WHT: %v, Expected Status Code: %d, Got: %d, Error: %v",
					tt.name, tt.totalIncome, tt.wht, tt.statusCode, rec.Code, err)
				return
			}

			fmt.Printf(color.GreenString("%s passed. Income: %v, WHT: %v\n"), tt.name, tt.totalIncome, tt.wht)
		})
	}
}
