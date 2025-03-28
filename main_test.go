package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	entity "github.com/takassh/super-shiharai-kun/entity"
	repo "github.com/takassh/super-shiharai-kun/repository"
	req "github.com/takassh/super-shiharai-kun/request"
)

func TestMain(m *testing.M) {
	repo.PopulateCompany(companyRepo)
	repo.PopulateUser(userRepo)
	repo.PopulateClient(companyRepo)
	repo.PopulateBankAccount(bankRepo)
	repo.PopulateInvoice(invoiceRepo)
	loadErrors()

	status := m.Run()

	os.Exit(status)
}

func TestCreateInvoices(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	request := req.CreateInvoice{
		Amount:   10000,
		ClientID: 2,
		DueDate:  "2025-03-30",
	}
	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/api/invoices", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &jwt.Token{
		Claims: &Claims{
			UserID: "1",
		},
	})
	// Assertions
	if assert.NoError(t, CreateInvoice(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		// parse
		var resp entity.Invoice
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		// assert
		assert.Equal(t, uint(3), resp.ID)
		assert.Equal(t, uint(1), resp.CompanyID)
		assert.Equal(t, uint(2), resp.ClientID)
		assert.NotEqual(t, "", resp.IssuedAt)
		assert.Equal(t, float64(10000), resp.Amount)
		assert.Equal(t, 400.0, resp.Fee)
		assert.Equal(t, 0.04, resp.FeeRate)
		assert.Equal(t, 40.0, resp.Tax)
		assert.Equal(t, 0.1, resp.TaxRate)
		assert.Equal(t, 10440.0, resp.TotalAmount)
		assert.Equal(t, "2025-03-30T00:00:00Z", resp.DueDate.Format("2006-01-02T15:04:05Z"))
		assert.Equal(t, entity.StatusUnprocessed, resp.Status)
	}
}

func TestGetInvoices(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodGet, "/api/invoices?start_date=2025-04-01&due_date=2025-04-30", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &jwt.Token{
		Claims: &Claims{
			UserID:    "1",
			CompanyID: "1",
		},
	})
	// Assertions
	if assert.NoError(t, GetInvoices(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		// parse
		var resp []entity.Invoice
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		// assert
		assert.Len(t, resp, 1)
		r := resp[0]
		assert.Equal(t, uint(1), r.ID)
		assert.Equal(t, uint(1), r.CompanyID)
		assert.Equal(t, uint(1), r.ClientID)
		assert.Equal(t, "2025-03-28T00:00:00Z", r.IssuedAt.Format("2006-01-02T15:04:05Z"))
		assert.Equal(t, float64(10000), r.Amount)
		assert.Equal(t, 400.0, r.Fee)
		assert.Equal(t, 0.04, r.FeeRate)
		assert.Equal(t, 40.0, r.Tax)
		assert.Equal(t, 0.1, r.TaxRate)
		assert.Equal(t, 10440.0, r.TotalAmount)
		assert.Equal(t, "2025-04-28T00:00:00Z", r.DueDate.Format("2006-01-02T15:04:05Z"))
		assert.Equal(t, entity.StatusUnprocessed, r.Status)
	}
}
