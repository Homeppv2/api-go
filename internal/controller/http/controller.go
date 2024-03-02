package http

import (
	"api-go/internal/controller"
	"api-go/internal/controller/middleware"
	"api-go/internal/entity"
	"api-go/pkg/errors"

	"github.com/gofiber/fiber/v2"
)

type ControllerHandler struct {
	controllerService controller.ControllerService
	mdwManager        middleware.MdwManager
}

func NewControllerHandler(
	controllerService controller.ControllerService,
	mdwManager middleware.MdwManager,
) *ControllerHandler {
	return &ControllerHandler{
		controllerService: controllerService,
		mdwManager:        mdwManager,
	}
}

func (h *ControllerHandler) create() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var p entity.CreateControllerRequest

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
		res, err := h.controllerService.Create(c.Context(), p.HwKey)
		if err != nil {
			return HandleError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(newResp(
			fiber.Map{
				"controller": res,
			},
			fiber.Map{
				"success": true,
			},
		))
	}
}

func (h *ControllerHandler) getByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var p entity.GetControllerByIDRequest

		if err := c.ParamsParser(&p); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				newErrResp(err),
			)
		}

		if err := p.Validate(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				newErrResp(errors.ErrValidationError),
			)
		}
		res, err := h.controllerService.GetByID(c.Context(), p.ID)
		if err != nil {
			return HandleError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(newResp(
			fiber.Map{
				"controller": res,
			},
			fiber.Map{
				"success": true,
			},
		))
	}
}

func (h *ControllerHandler) getByHwKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var p entity.GetControllerByHwKeyRequest

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
		res, err := h.controllerService.GetByHwKey(c.Context(), p.HwKey)
		if err != nil {
			return HandleError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(newResp(
			fiber.Map{
				"controller": res,
			},
			fiber.Map{
				"success": true,
			},
		))
	}
}

func (h *ControllerHandler) getByIsUsedBy() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var p entity.GetControllerByIsUsedByRequest

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
		res, err := h.controllerService.GetByIsUsedBy(c.Context(), p.IsUsedBy)
		if err != nil {
			return HandleError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(newResp(
			fiber.Map{
				"controller": res,
			},
			fiber.Map{
				"success": true,
			},
		))
	}
}

func (h *ControllerHandler) updateIsUsed() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var p entity.UpdateControllerIsUsedByRequest

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
		res, err := h.controllerService.UpdateIsUsed(c.Context(), p)
		if err != nil {
			return HandleError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(newResp(
			fiber.Map{
				"controller": res,
			},
			fiber.Map{
				"success": true,
			},
		))
	}
}

func (h *ControllerHandler) delete() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var p entity.DeleteControllerByIDRequest

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
		err := h.controllerService.Delete(c.Context(), p.ID)
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

func (h *ControllerHandler) Register(r fiber.Router) {
	r.Post("create",
		// h.mdwManager.SessionValidate(),
		h.create())

	r.Get("id/:id",
		h.mdwManager.SessionValidate(),
		h.getByID())

	r.Get("hw-key/:hwKey",
		h.mdwManager.SessionValidate(),
		h.getByHwKey())

	r.Get("is-used-by/:isUsedBy",
		h.mdwManager.SessionValidate(),
		h.getByIsUsedBy())

	r.Patch("is-used-by/",
		h.mdwManager.SessionValidate(),
		h.updateIsUsed())

	r.Delete(":id",
		h.mdwManager.SessionValidate(),
		h.delete())
}
