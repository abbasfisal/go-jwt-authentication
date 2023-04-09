package models

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

func (u *User) GeneratePassword() error {
	bytePassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	u.Password = string(bytePassword)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
func (u *User) Validate() bool {
	fmt.Println(u)
	if len(u.Email) == 0 || len(u.Email) <= 3 || len(u.Password) == 0 || len(u.Password) < 4 {
		return false
	}
	return true
}
