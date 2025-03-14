package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
)

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

var (
	clientId     = ""
	clientSecret = ""
	redirectUri  = "http://localhost:4000/api/auth/callback"
	tokenUrl     = "https://oauth2.googleapis.com/token"
	userInfoUrl  = "https://www.googleapis.com/oauth2/v1/userinfo"
)

var client = &http.Client{}

func GetAccessTokenByAuthorizationCode(code string) (*tokenResponse, error) {
	resp, err := client.PostForm(tokenUrl, map[string][]string{
		"client_id":     {clientId},
		"client_secret": {clientSecret},
		"code":          {code},
		"redirect_uri":  {redirectUri},
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
	req, err := http.NewRequest("GET", userInfoUrl, nil)
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
