package models

import "time"

type Client struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"not null"`
	LastName  string    `json:"last_name" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	BirthDay  time.Time `json:"birth_day" gorm:"not null"`
	Age       int       `json:"age" gorm:"not null"`
	Telephone string    `json:"telephone"`
}
