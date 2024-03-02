package http

import (
	"api-go/internal/controller"
	"api-go/internal/controller/middleware"
	"api-go/internal/entity"
	"api-go/pkg/errors"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService controller.UserService
	mdwManager  middleware.MdwManager
}

func NewUserHandler(
	userService controller.UserService,
	mdwManager middleware.MdwManager,
) *UserHandler {
	return &UserHandler{
		userService: userService,
		mdwManager:  mdwManager,
	}
}

func (h *UserHandler) register() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var p entity.RegisterUserRequest

		if err := c.BodyParser(&p); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				newErrResp(errors.ErrValidationError),
			)
		}

		if err := p.Validate(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				newErrResp(err),
			)
		}
		res, err := h.userService.Register(c.Context(), p)
		if err != nil {
			return HandleError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(newResp(
			fiber.Map{
				"user": res,
			},
			fiber.Map{
				"success": true,
			},
		))
	}
}

func (h *UserHandler) getByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var p entity.GetUserByIDRequest

		if err := c.ParamsParser(&p); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				newErrResp(errors.ErrValidationError),
			)
		}

		if err := p.Validate(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				newErrResp(err),
			)
		}
		res, err := h.userService.GetByID(c.Context(), p.ID)
		if err != nil {
			return HandleError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(newResp(
			fiber.Map{
				"user": res,
			},
			fiber.Map{
				"success": true,
			},
		))
	}
}

func (h *UserHandler) getByEmail() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var p entity.GetUserByEmailRequest

		if err := c.ParamsParser(&p); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				newErrResp(errors.ErrValidationError),
			)
		}

		if err := p.Validate(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				newErrResp(err),
			)
		}
		res, err := h.userService.GetByEmail(c.Context(), p.Email)
		if err != nil {
			return HandleError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(newResp(
			fiber.Map{
				"user": res,
			},
			fiber.Map{
				"success": true,
			},
		))
	}
}

func (h *UserHandler) Register(r fiber.Router) {
	r.Post("register",
		h.register())

	r.Get("id/:id",
		h.mdwManager.SessionValidate(),
		h.getByID())

	r.Get("email/:email",
		h.mdwManager.SessionValidate(),
		h.getByEmail())
}
