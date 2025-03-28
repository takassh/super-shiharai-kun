package entity

import (
	"errors"
	"time"
)

type Status string

const (
	StatusUnprocessed Status = "未処理"
	StatusProcessing  Status = "処理中"
	StatusPaid        Status = "支払い済み" // final status
	StatusError       Status = "エラー"   // final status
)

type Invoice struct {
	ID          uint `gorm:"primaryKey"`
	CompanyID   uint
	ClientID    uint
	IssuedAt    time.Time
	Amount      float64
	Fee         float64
	FeeRate     float64
	Tax         float64
	TaxRate     float64
	TotalAmount float64
	DueDate     time.Time
	Status      Status
}

func (i *Invoice) CalculateFee() error {
	if i.Amount == 0 {
		return errors.New("amount must be greater than 0")
	}
	if i.FeeRate == 0 {
		return errors.New("fee rate must be greater than 0")
	}
	if i.TaxRate == 0 {
		return errors.New("tax rate must be greater than 0")
	}

	i.Fee = i.Amount * i.FeeRate

	return nil
}

func (i *Invoice) CalculateFeeRate() error {
	if i.Amount == 0 {
		return errors.New("amount must be greater than 0")
	}
	if i.Fee == 0 {
		return errors.New("fee must be greater than 0")
	}
	i.FeeRate = i.Fee / i.Amount
	return nil
}

func (i *Invoice) CalculateTaxRate() error {
	if i.Fee == 0 {
		return errors.New("fee must be greater than 0")
	}
	if i.Tax == 0 {
		return errors.New("tax must be greater than 0")
	}
	i.TaxRate = i.Tax / i.Fee
	return nil
}

func (i *Invoice) CalculateTax() error {
	if i.Fee == 0 {
		return errors.New("fee must be greater than 0")
	}
	if i.TaxRate == 0 {
		return errors.New("tax rate must be greater than 0")
	}

	i.Tax = i.Fee * i.TaxRate
	return nil
}

func (i *Invoice) CalculateTotalAmount() error {
	if i.Amount == 0 {
		return errors.New("amount must be greater than 0")
	}
	if i.FeeRate == 0 {
		return errors.New("fee rate must be greater than 0")
	}
	if i.TaxRate == 0 {
		return errors.New("tax rate must be greater than 0")
	}
	i.TotalAmount = i.Amount + i.Fee + i.Tax

	return nil
}
