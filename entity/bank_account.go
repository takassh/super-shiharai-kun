package entity

type BankAccount struct {
	ID          uint `gorm:"primaryKey"`
	CompanyID   uint
	BankName    string
	Branch      string
	AccountNo   string
	AccountName string
}
