package repository

import (
	e "github.com/takassh/super-shiharai-kun/entity"
	"gorm.io/gorm"
)

var userRepositoryEntities = []any{
	&e.User{},
}

type UserRepository interface {
	Save(e.User) error
	FindByID(uint) (e.User, error)
}

func (r *RepositoryFactory) NewUserRepository() UserRepository {
	if r.autoMigrate {
		r.connect().AutoMigrate(userRepositoryEntities...)
	}

	return &UserRepositoryImpl{r.connect()}
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func (h *UserRepositoryImpl) Save(user e.User) error {
	result := h.db.Save(&user)
	return result.Error
}

func (h *UserRepositoryImpl) FindByID(id uint) (e.User, error) {
	var user e.User
	result := h.db.First(&user, id)
	return user, result.Error
}
