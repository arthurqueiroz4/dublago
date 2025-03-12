package auth

import (
	"errors"
	"fmt"
	"tradutor-dos-crias/config"
	"tradutor-dos-crias/singleton"
	"tradutor-dos-crias/user"

	"github.com/gofiber/fiber/v3"
)

func Callback(c fiber.Ctx) error {
	authCode := c.Query("code")
	if authCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(struct{ Message string }{
			Message: "It's not possible get authorization code",
		})
	}

	tokenResponse, err := GetAccessTokenByAuthorizationCode(authCode)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	userInfo, err := GetUserInfoByAccessToken(tokenResponse.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	u := &user.User{
		Email: userInfo["email"].(string),
		Name:  userInfo["name"].(string),
		SsoId: userInfo["id"].(string),
	}

	err = singleton.UserService.Create(u)
	if err != nil && !errors.Is(err, user.ErrUserAlreadyExists) {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	c.Response().Header.Add("Location", fmt.Sprintf("%s?access_token=%s", config.FrontendRedirectUrl, tokenResponse.AccessToken))
	c.Status(301)
	return nil
}
