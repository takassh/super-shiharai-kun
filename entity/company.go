package entity

import "gorm.io/gorm"

type Company struct {
	ID         uint
	Name       string
	CEO        string
	Phone      string
	PostalCode string
	Address    string
	Users      []User
	Clients    []Company `gorm:"many2many:company_clients;"` // from companies to clients
	Companies  []Company `gorm:"many2many:company_clients;"` // from clients to companies
}

type CompanyClient struct {
	gorm.Model
	CompanyID uint `gorm:"index" json:"company_id"`
	ClientID  uint `gorm:"index" json:"client_id"`
}
