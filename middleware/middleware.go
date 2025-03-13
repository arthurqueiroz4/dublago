package middleware

import (
	"log/slog"
	"tradutor-dos-crias/auth"

	"github.com/gofiber/fiber/v3"
)

func Authentication(c fiber.Ctx) error {
	slog.Info("Reached here")

	accessToken := c.Get("Authorization")
	userInfo, err := auth.GetUserInfoByAccessToken(accessToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "failed to get user info"})
	}

	c.Locals("user.name", userInfo["name"])
	c.Locals("user.email", userInfo["email"])
	c.Locals("user.sso_id", userInfo["id"])

	return c.Next()
}

func DefaultAuthorization(c fiber.Ctx) error {
	email := c.Locals("user.email").(string)

	return c.Next()
}
