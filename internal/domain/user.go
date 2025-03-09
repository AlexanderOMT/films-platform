package domain

type User struct {
	Id       uint   `gorm:"primaryKey"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
