package domain

type User struct {
	Id       uint `gorm:"primaryKey"`
	Username string
	Password string
}
