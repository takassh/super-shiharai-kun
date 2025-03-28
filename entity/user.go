package entity

type User struct {
	ID        uint `gorm:"primaryKey"`
	CompanyID uint
	Name      string
	Email     string
	Password  string
}
