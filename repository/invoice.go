package repository

import (
	e "github.com/takassh/super-shiharai-kun/entity"
	"gorm.io/gorm"
)

var invoiceRepositoryEntities = []any{
	&e.Invoice{},
}

type InvoiceRepository interface {
	Save(*e.Invoice) error
	Find() ([]e.Invoice, error)
	FindByDate(uint, string, string) ([]e.Invoice, error)
}

func (r *RepositoryFactory) NewInvoiceRepository() InvoiceRepository {
	if r.autoMigrate {
		r.connect().AutoMigrate(invoiceRepositoryEntities...)
	}

	return &InvoiceRepositoryImpl{r.connect()}
}

type InvoiceRepositoryImpl struct {
	db *gorm.DB
}

func (h *InvoiceRepositoryImpl) Save(invoice *e.Invoice) error {
	result := h.db.Save(&invoice)
	return result.Error
}

func (h *InvoiceRepositoryImpl) Find() ([]e.Invoice, error) {
	var invoices []e.Invoice
	result := h.db.Order("date DESC").Find(&invoices)
	return invoices, result.Error
}

// FindByDate returns invoices by date in descending order
func (h *InvoiceRepositoryImpl) FindByDate(companyID uint, startDate string, endDate string) ([]e.Invoice, error) {
	var invoices []e.Invoice
	query := h.db.Order("due_date DESC").Where("company_id = ?", companyID)
	if startDate != "" {
		query = query.Where("due_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("due_date <= ?", endDate)
	}
	result := query.Find(&invoices)
	return invoices, result.Error
}
