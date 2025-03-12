package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	ClientId            string
	ClientSecret        string
	RedirectUri         string
	TokenUrl            string
	UserInfoUrl         string
	FrontendRedirectUrl string
	DatabaseUrl         string
)

func LoadEnvironment() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error on load environment", err)
	}

	ClientId = os.Getenv("CLIENT_ID")
	ClientSecret = os.Getenv("CLIENT_SECRET")
	RedirectUri = os.Getenv("REDIRECT_URI")
	TokenUrl = os.Getenv("TOKEN_URL")
	UserInfoUrl = os.Getenv("USER_INFO_URL")
	FrontendRedirectUrl = os.Getenv("FRONTEND_REDIRECT_URL")
	DatabaseUrl = os.Getenv("DATABASE_URL")

	if ClientId == "" {
		log.Fatal("ClientId is empty")
	}
	if ClientSecret == "" {
		log.Fatal("ClientSecret is empty")
	}
	if FrontendRedirectUrl == "" {
		log.Fatal("FrontendRedirect is empty")
	}
	if DatabaseUrl == "" {
		log.Fatal("DatabaseUrl is empty")
	}
	if RedirectUri == "" {
		RedirectUri = "http://localhost:4000/api/auth/callback"
	}
	if TokenUrl == "" {
		TokenUrl = "https://oauth2.googleapis.com/token"
	}
	if UserInfoUrl == "" {
		UserInfoUrl = "https://www.googleapis.com/oauth2/v1/userinfo"
	}
}
