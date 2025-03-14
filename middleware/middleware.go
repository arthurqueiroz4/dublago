package middleware

import (
	"tradutor-dos-crias/auth"
	"tradutor-dos-crias/singleton"

	"github.com/gofiber/fiber/v3"
)

func Authentication(c fiber.Ctx) error {
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

	user, err := singleton.UserService.GetByEmail(email)
	if err != nil || user.ID == 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "user not registred"})
	}

	return c.Next()
}
