package repository

import (
	e "github.com/takassh/super-shiharai-kun/entity"
	"gorm.io/gorm"
)

var companyRepositoryEntities = []any{
	&e.Company{},
}

type CompanyRepository interface {
	CreateAll([]e.Company) error
	Save(e.Company) error
	FindByID(uint) (e.Company, error)
	FindByDate(string) ([]e.Company, error)
	AppendClient(e.Company, e.Company) error
	RemoveClient(e.Company, e.Company) error
}

func (r *RepositoryFactory) NewCompanyRepository() CompanyRepository {
	if r.autoMigrate {
		r.connect().AutoMigrate(companyRepositoryEntities...)
	}

	return &CompanyRepositoryImpl{r.connect()}
}

type CompanyRepositoryImpl struct {
	db *gorm.DB
}

func (h *CompanyRepositoryImpl) CreateAll(companies []e.Company) error {
	result := h.db.Create(companies)
	return result.Error
}

func (h *CompanyRepositoryImpl) Save(company e.Company) error {
	result := h.db.Save(&company)
	return result.Error
}

func (h *CompanyRepositoryImpl) AppendClient(company e.Company, client e.Company) error {
	return h.db.Model(&company).Association("Clients").Append(&client)
}

func (h *CompanyRepositoryImpl) RemoveClient(company e.Company, client e.Company) error {
	return h.db.Model(&company).Association("Clients").Delete(&client)
}

func (h *CompanyRepositoryImpl) AppendUser(company e.Company, user e.User) error {
	return h.db.Model(&company).Association("Users").Append(&user)
}

func (h *CompanyRepositoryImpl) RemoveUser(company e.Company, user e.User) error {
	return h.db.Model(&company).Association("Users").Delete(&user)
}

func (h *CompanyRepositoryImpl) AppendBankAccount(company e.Company, bankAccount e.BankAccount) error {
	return h.db.Model(&company).Association("BankAccounts").Append(&bankAccount)
}

func (h *CompanyRepositoryImpl) RemoveBankAccount(company e.Company, bankAccount e.BankAccount) error {
	return h.db.Model(&company).Association("BankAccounts").Delete(&bankAccount)
}

func (h *CompanyRepositoryImpl) FindByID(id uint) (e.Company, error) {
	var company e.Company
	result := h.db.First(&company, id)
	return company, result.Error
}

// FindByDate returns companies by date in descending order
func (h *CompanyRepositoryImpl) FindByDate(yearAndMonth string) ([]e.Company, error) {
	var companies []e.Company
	result := h.db.Where("date LIKE ?", yearAndMonth+"%").Order("date DESC").Find(&companies)
	return companies, result.Error
}
