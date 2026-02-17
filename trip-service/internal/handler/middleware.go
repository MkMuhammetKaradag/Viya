package handler

import (
	"errors"
	"trip-service/internal/database"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var validate = validator.New()

func HandleBasic[R Request, Res Response](handler BasicHandler[R, Res]) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req R

		if err := parseRequest(c, &req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := validate.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "validation failed", "details": err.Error()})
		}

		ctx := c.Context()
		res, err := handler.Handle(ctx, &req)

		if err != nil {
			status := getStatusCodeFromError(err)
			return c.Status(status).JSON(fiber.Map{"error": err.Error()})
		}
		if c.Method() == fiber.MethodPost {
			return c.Status(fiber.StatusCreated).JSON(res)
		}
		return c.JSON(res)
	}
}
func HandleWithFiber[R Request, Res Response](handler FiberHandler[R, Res]) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req R

		if err := parseRequest(c, &req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := validate.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "validation failed", "details": err.Error()})
		}

		res, err := handler.Handle(c, &req)

		if err != nil {
			status := getStatusCodeFromError(err)
			return c.Status(status).JSON(fiber.Map{"error": err.Error()})
		}
		// if c.Method() == fiber.MethodPost {
		// 	return c.Status(fiber.StatusCreated).JSON(res)
		// }
		return c.JSON(res)
	}
}
func parseRequest[R any](c fiber.Ctx, req *R) error {

	if err := c.Bind().Body(req); err != nil {
		return err
	}
	if err := c.Bind().Query(req); err != nil {
		return err
	}
	if err := c.Bind().URI(req); err != nil {
		return err
	}
	if err := c.Bind().Header(req); err != nil {
		return err
	}
	return nil
}
func getStatusCodeFromError(err error) int {
	switch {

	case errors.Is(err, database.ErrDuplicateResource):
		return fiber.StatusConflict

	default:
		return fiber.StatusInternalServerError
	}
}
