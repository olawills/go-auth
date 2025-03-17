package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username" gorm:"unique"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"-" gorm:"column:password" validate:"required,min=6,max=20"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Login struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Register struct {
	FirstName string `json:"first_name" binding:"required" validate:"required,min=2,max=20"`
	LastName  string `json:"last_name" binding:"required" validate:"required,min=2,max=20"`
	Username  string `json:"username" binding:"required" gorm:"unique"`
	Email     string `json:"email" binding:"required" validate:"required,email"`
	Password  string `json:"password" binding:"required" validate:"required,min=6,max=20"`
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		// panic(err)
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	// panic(err)
	if err != nil {

		return err
	}
	return nil
}
