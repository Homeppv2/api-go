package entity

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"gopkg.in/guregu/null.v4"
)

type Controller struct {
	ID       int64    `json:"id" db:"id"`
	HwKey    string   `json:"hwKey" db:"hw_key"`
	IsUsed   *bool    `json:"isUsed" db:"is_used"`
	IsUsedBy null.Int `json:"isUsedBy" db:"is_used_by"`
}

type CreateControllerRequest struct {
	HwKey string `json:"hwKey" db:"hw_key"`
}

func (s CreateControllerRequest) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.HwKey, validation.Required.Error("ключ контроллера обязателен для заполнения"),
			validation.Length(3, 3).Error("ключ контроллера должен быть длиной 3 символа")),
	)
}

type CreateControllerDTO struct {
	HwKey    string
	IsUsed   bool
	IsUsedBy int64
}

type GetControllerByIDRequest struct {
	ID int64 `params:"id"`
}

func (s GetControllerByIDRequest) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.ID, validation.Required.Error("идентификатор контроллера обязателен для заполнения"),
			validation.Min(1).Error("идентификатор контроллера должен быть положительным целым числом")),
	)
}

type GetControllerByHwKeyRequest struct {
	HwKey string `params:"hwKey"`
}

func (s GetControllerByHwKeyRequest) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.HwKey, validation.Required.Error("ключ контроллера обязателен для заполнения"),
			validation.Length(3, 3).Error("ключ контроллера должен быть длиной 3 символа")),
	)
}

type GetControllerByIsUsedByRequest struct {
	IsUsedBy int64 `params:"isUsedBy"`
}

func (s GetControllerByIsUsedByRequest) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.IsUsedBy, validation.Required.Error("идентификатор пользователя обязателен для заполнения"),
			validation.Min(1).Error("идентификатор пользователя должен быть положительным целым числом")))
}

type UpdateControllerIsUsedByRequest struct {
	ID       int64
	IsUsed   *bool
	IsUsedBy null.Int
}

func (s UpdateControllerIsUsedByRequest) Validate() error {
	if s.ID <= 0 {
		return errors.New("идентификатор контроллера обязателен для заполнения")
	}

	if s.IsUsed == nil {
		return errors.New("флаг использования обязателен для заполнения")
	}

	if s.IsUsedBy.Valid {
		if s.IsUsedBy.ValueOrZero() <= 0 {
			return errors.New("идентификатор пользователя должен быть положительным целым числом")
		}
		if !*s.IsUsed {
			return errors.New("флаг использования c неправильным значением с данным значением идентификатора пользователя")
		}
	} else {
		if *s.IsUsed {
			return errors.New("идентификатор пользователя обязателен для заполнения с данным флагом использования")
		}
	}

	return nil
}

type DeleteControllerByIDRequest struct {
	ID int64 `params:"id"`
}

func (s DeleteControllerByIDRequest) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.ID, validation.Required.Error("идентификатор контроллера обязателен для заполнения"),
			validation.Min(1).Error("идентификатор контроллера должен быть положительным целым числом")),
	)
}
