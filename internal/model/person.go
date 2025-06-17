package model

import "time"

type Person struct {
	ID         uint      `gorm:"primary_key;" json:"id"`
	Name       string    `json:"name"`
	Surname    string    `json:"surname"`
	Patronymic string    `json:"patronymic"`
	Age        int       `json:"age"`
	Gender     string    `json:"gender"`
	National   string    `json:"national"`
	CreatedAt  time.Time `json:"created_at"`
}

type InputPerson struct {
	Name       string `json:"name" binding:"required"`
	Surname    string `json:"surname" binding:"required"`
	Patronymic string `json:"patronymic,omitempty"`
}
