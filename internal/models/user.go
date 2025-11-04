package models

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Password string `json:"password"` // Password hash
}

// Хеширование пароля(используется bcrypt)
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Сравнение пароля с хешем из БД
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) Validate() error {
	if strings.TrimSpace(u.UserName) == "" {
		return fmt.Errorf("Имя пользователя не может быть пустым")
	}

	if len(u.UserName) < 2 {
		return fmt.Errorf("Имя пользователя должно содержать минимум 2 символа")
	}

	if len(u.UserName) > 255 {
		return fmt.Errorf("Имя пользователя не должно превышать 255 символов")
	}

	if strings.TrimSpace(u.Password) == "" {
		return fmt.Errorf("Пароль не может быть пустым")
	}

	return nil
}
