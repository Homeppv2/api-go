package middleware

import (
	"github.com/gofiber/fiber/v2"
)

type MdwManager struct {
	authService AuthService
}

func NewMiddlewareManager(authService AuthService) *MdwManager {
	return &MdwManager{
		authService: authService,
	}
}

func (m *MdwManager) SessionValidate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get(fiber.HeaderAuthorization)
		session, err := m.authService.CheckSession(c.Context(), token)
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		c.Locals("userID", session.UserID)
		c.Locals("token", session.Token)

		return c.Next()
	}
}
