package request

import (
	"errors"
	"time"

	e "github.com/takassh/super-shiharai-kun/entity"
	u "github.com/takassh/super-shiharai-kun/util"
)

type CreateInvoice struct {
	ClientID uint    `json:"client_id" validate:"required"` // frontend should provide client_id instead of client_name
	Amount   float64 `json:"amount" validate:"required"`
	DueDate  string  `json:"due_date" validate:"required"` // frontend should provide due_date as UTC
}

func (c *CreateInvoice) ValidateValue() error {
	if c.ClientID == 0 {
		return errors.New("client_id is not provided")
	}
	if c.Amount == 0 {
		return errors.New("amount is not provided")
	}
	dueDate, err := u.Atot(c.DueDate)
	if err != nil {
		return errors.New("due_date is not provided")
	}

	now := time.Now().UTC()
	if dueDate.Before(now) {
		return errors.New("due_date is past")
	}

	return nil
}

// ToEntity converts CreateInvoice to Invoice entity
func (c *CreateInvoice) ToEntity(companyID uint, status e.Status, feeRate float64, taxRate float64) (*e.Invoice, error) {
	err := c.ValidateValue()
	if err != nil {
		return nil, err
	}

	dueDate, _ := u.Atot(c.DueDate)
	entity := &e.Invoice{
		CompanyID: companyID,
		ClientID:  c.ClientID,
		IssuedAt:  time.Now().UTC(),
		Amount:    c.Amount,
		FeeRate:   feeRate,
		TaxRate:   taxRate,
		DueDate:   dueDate,
		Status:    status,
	}

	entity.CalculateFee()
	entity.CalculateTax()
	entity.CalculateTotalAmount()
	return entity, nil
}
