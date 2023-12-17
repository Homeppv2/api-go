package entity

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type User struct {
	ID           int64  `json:"id" db:"id"`
	Username     string `json:"username" db:"username"`
	Email        string `json:"email" db:"email"`
	HashPassword string `json:"hashPassword" db:"hash_password"`
}

type RegisterUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	HwKey    string `json:"hwKey"`
}

func (s RegisterUserRequest) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Username, validation.Required.Error("username пользователя обязателен для заполнения"),
			validation.Length(6, 20).Error("username пользователя должен быть длиной от 6 до 20 символов")),
		validation.Field(&s.Email, validation.Required.Error("email пользователя обязателен для заполнения"),
			is.Email.Error("неправильный формат email пользователя")),
		validation.Field(&s.Password, validation.Required.Error("пароль пользователя обязателен для заполнения"),
			validation.Length(6, 30).Error("пароль пользователя должен быть длиной от 6 до 30 символов")),
		validation.Field(&s.HwKey, validation.Required.Error("ключ контроллера обязателен для заполнения"),
			validation.Length(3, 3).Error("ключ контроллера должен быть длиной 3 символа")),
	)
}

type CreateUserDTO struct {
	Username     string
	Email        string
	HashPassword string
}

type GetUserByIDRequest struct {
	ID int64 `params:"id"`
}

func (s GetUserByIDRequest) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.ID, validation.Required.Error("идентификатор пользователя обязателен для заполнения"),
			validation.Min(1).Error("идентификатор пользователя должен быть положительным целым числом")),
	)
}

type GetUserByEmailRequest struct {
	Email string `params:"email"`
}

func (s GetUserByEmailRequest) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Email, validation.Required.Error("email пользователя обязателен для заполнения"),
			is.Email.Error("неправильный формат email пользователя")),
	)
}
