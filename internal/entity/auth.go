package entity

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Session struct {
	Token  string
	UserID int64
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s LoginRequest) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Email, validation.Required.Error("email пользователя обязателен для заполнения"),
			is.Email.Error("неправильный формат email пользователя")),
		validation.Field(&s.Password, validation.Required.Error("пароль пользователя обязателен для заполнения")),
	)
}

type LoginResponse struct {
	User    User    `json:"user"`
	Session Session `json:"session"`
}
