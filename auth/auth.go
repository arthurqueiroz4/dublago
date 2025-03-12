package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

var (
	clientId         = "343946256424-1qh1f475m4eofq5i4g0m5tk38pnp1h0l.apps.googleusercontent.com"
	clientSecret     = ""
	redirectUri      = "http://localhost:4000/api/auth/callback"
	frontRedirectUrl = "http://localhost:3000/dashboard"
	tokenUrl         = "https://oauth2.googleapis.com/token"
	userInfoUrl      = "https://www.googleapis.com/oauth2/v1/userinfo"
)

func Callback(c fiber.Ctx) error {
	authCode := c.Query("code")
	if authCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(struct{ Message string }{
			Message: "It's not possible get authorization code",
		})
	}

	resp, err := http.PostForm(tokenUrl, map[string][]string{
		"client_id":     {clientId},
		"client_secret": {clientSecret},
		"code":          {authCode},
		"redirect_uri":  {redirectUri},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao trocar código por token",
		})
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":             "Erro ao comunicar-se com API do Google",
			"messageFromGoogle": string(body),
		})
	}

	slog.Info("Response from /token",
		"body", string(body))

	var tokenResponse struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&tokenResponse); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao decodificar resposta do Google",
		})
	}

	req, err := http.NewRequest("GET", userInfoUrl, nil)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Erro na criação da request para API de userinfo"})
	}

	req.Header.Add("Authorization", "Bearer "+tokenResponse.AccessToken)

	client := &http.Client{}

	resp, err = client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Erro na chamada da request para API de userinfo"})
	}
	defer resp.Body.Close()

	var userInfo map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Erro ao decodificar request para API de userinfo"})
	}

	slog.Info("Authentication completed",
		"username", userInfo["name"],
		"email", userInfo["email"],
		"sso_id", userInfo["id"])

	c.Response().Header.Add("Location", frontRedirectUrl)
	c.Status(301)
	return nil
}
