package model

type User struct {
	ID   string `gorm:"primaryKey"`
	Name string
}
