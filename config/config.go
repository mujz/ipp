package config

import (
	"os"
	"strconv"
	"time"
)

var (
	// Server vars
	Port    = os.Getenv("PORT")
	BaseURL = os.Getenv("BASE_URL")
	WebURL  = os.Getenv("WEB_URL")

	// Database vars
	DBName     = os.Getenv("DB_NAME")
	DBUser     = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBHost     = os.Getenv("DB_HOST")
	DBPort     = os.Getenv("DB_PORT")
	DBSSLMode  = os.Getenv("DB_SSL_MODE")

	// Facebook vars
	FBAppID     = os.Getenv("FB_APP_ID")
	FBAppSecret = os.Getenv("FB_APP_SECRET")

	// Authentication vars
	AuthSecretKey               = []byte(os.Getenv("SECRET_KEY"))
	AuthTokenExpirationInterval time.Duration
)

func init() {
	if Port == "" {
		Port = "8000"
	}
	if BaseURL == "" {
		BaseURL = "http://localhost"
	}
	if BaseURL == "" {
		BaseURL = "http://localhost:8080"
	}

	if DBName == "" {
		DBName = "ipp"
	}
	if DBUser == "" {
		DBUser = "mujtaba"
	}
	if DBPassword == "" {
		DBPassword = "thinkific"
	}
	if DBHost == "" {
		DBHost = "localhost"
	}
	if DBPort == "" {
		DBPort = "5432"
	}
	if DBSSLMode == "" {
		DBSSLMode = "disable"
	}

	if numOfSeconds, err := strconv.Atoi(
		os.Getenv("AUTH_TOKEN_EXPIRATION_INTERVAL_IN_SECONDS"),
	); err != nil {
		AuthTokenExpirationInterval = 3600 * time.Second
	} else {
		AuthTokenExpirationInterval = time.Duration(numOfSeconds) * time.Second
	}

	if AuthSecretKey == nil {
		AuthSecretKey = []byte("ipp,secret")
	}
}
