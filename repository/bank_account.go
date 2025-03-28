package repository

import (
	e "github.com/takassh/super-shiharai-kun/entity"
	"gorm.io/gorm"
)

var bankAccountRepositoryEntities = []any{
	&e.BankAccount{},
}

type BankAccountRepository interface {
	Save(e.BankAccount) error
	FindByID(uint) (e.BankAccount, error)
}

func (r *RepositoryFactory) NewBankAccountRepository() BankAccountRepository {
	if r.autoMigrate {
		r.connect().AutoMigrate(bankAccountRepositoryEntities...)
	}

	return &BankAccountRepositoryImpl{r.connect()}
}

type BankAccountRepositoryImpl struct {
	db *gorm.DB
}

func (h *BankAccountRepositoryImpl) Save(bankAccount e.BankAccount) error {
	result := h.db.Save(&bankAccount)
	return result.Error
}

func (h *BankAccountRepositoryImpl) FindByID(id uint) (e.BankAccount, error) {
	var bankAccount e.BankAccount
	result := h.db.First(&bankAccount, id)
	return bankAccount, result.Error
}
