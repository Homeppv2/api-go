package http

import (
	"errors"

	errorsPkg "api-go/pkg/errors"

	"github.com/gofiber/fiber/v2"
)

var errCodeMap = map[errorsPkg.ErrCode]int{
	errorsPkg.ErrCodeBadRequest:      fiber.StatusBadRequest,
	errorsPkg.ErrCodeUnknown:         fiber.StatusInternalServerError,
	errorsPkg.ErrCodeNotFound:        fiber.StatusNotFound,
	errorsPkg.ErrCodeInvalidArgument: fiber.StatusBadRequest,
}

func HandleError(ctx *fiber.Ctx, err error) error {
	appErr := &errorsPkg.Error{}
	if errors.As(err, &appErr) {
		c, ok := errCodeMap[appErr.Code()]
		if !ok {
			c = fiber.StatusInternalServerError
		}

		return ctx.Status(c).JSON(newErrResp(appErr))
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(newErrResp(err))
}
