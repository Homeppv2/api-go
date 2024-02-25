package http

import (
	"api-go/internal/controller"
	"api-go/internal/controller/middleware"
	"api-go/internal/entity"
	"api-go/pkg/errors"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService controller.AuthService
	mdwManager  middleware.MdwManager
}

func NewAuthHandler(
	authService controller.AuthService,
	mdwManager middleware.MdwManager,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		mdwManager:  mdwManager,
	}
}

func (h *AuthHandler) login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var p entity.LoginRequest

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
		res, err := h.authService.Login(c.Context(), p)
		if err != nil {
			return HandleError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(newResp(
			fiber.Map{
				"user":    res.User,
				"session": res.Session,
			},
			fiber.Map{
				"success": true,
			},
		))
	}
}

func (h *AuthHandler) logout() fiber.Handler {
	return func(c *fiber.Ctx) error {

		token := c.Locals("token").(string)

		err := h.authService.Logout(c.Context(), token)
		if err != nil {
			return HandleError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(newResp(
			fiber.Map{},
			fiber.Map{
				"success": true,
			},
		))
	}
}

func (h *AuthHandler) Register(r fiber.Router) {
	r.Post("login",
		h.login())

	r.Post("logout",
		h.mdwManager.SessionValidate(),
		h.logout())
}
