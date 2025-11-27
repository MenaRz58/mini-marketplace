package model

type User struct {
	ID       string `gorm:"primaryKey"`
	Name     string
	Email    string
	Password string
	IsAdmin  bool `gorm:"default:false"`
}
