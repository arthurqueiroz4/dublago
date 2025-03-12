package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"tradutor-dos-crias/config"
)

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

var client = &http.Client{}

func GetAccessTokenByAuthorizationCode(code string) (*tokenResponse, error) {
	resp, err := client.PostForm(config.TokenUrl, map[string][]string{
		"client_id":     {config.ClientId},
		"client_secret": {config.ClientSecret},
		"code":          {code},
		"redirect_uri":  {config.RedirectUri},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		return nil, errors.New("erro ao trocar código por token")
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, errors.New("erro ao comunicar-se com API do Google")
	}

	var tokenResponse tokenResponse

	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&tokenResponse); err != nil {
		return nil, errors.New("erro ao decodificar resposta do Google")
	}
	return &tokenResponse, nil
}

func GetUserInfoByAccessToken(accessToken string) (map[string]any, error) {
	req, err := http.NewRequest("GET", config.UserInfoUrl, nil)
	if err != nil {
		return nil, errors.New("erro na criação da request para API de userinfo")
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("erro na chamada da request para API de userinfo")
	}
	defer resp.Body.Close()

	var userInfo map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, errors.New("erro ao decodificar request para API de userinfo")
	}

	if userInfo["email"] == nil {
		return nil, errors.New("user not found by accessToken")
	}

	slog.Info("Authentication completed",
		"username", userInfo["name"],
		"email", userInfo["email"],
		"sso_id", userInfo["id"])

	return userInfo, nil
}
