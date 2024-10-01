package auth_utils

import (
	"api/config"
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/timewise-team/timewise-models/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"net/http"
	"time"
)

type GoogleOauthData struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	//Name          string `json:"name"`
	Picture string `json:"picture"`
}

var cfg, err = config.LoadConfig()

var GoogleOauth = oauth2.Config{
	ClientID:     cfg.GoogleOauth.ClientID,
	ClientSecret: cfg.GoogleOauth.ClientSecret,
	RedirectURL:  "",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

func VerifyGoogleToken(code string) ([]byte, error) {
	// compares the generated token string to the token retrieved from the parsed URL
	// converts authorization code into a token
	token, err := GoogleOauth.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, err
	}

	// this is done to prevent memory leakage
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// returns data of verified google user
	return data, nil
}

func GenerateJWTToken(user models.TwUser, secretKey string) (string, int, error) {
	// Định nghĩa thời gian hết hạn cho token (ví dụ: 2 giờ)
	expirationTime := time.Now().Add(2 * time.Hour).Unix()

	// Tạo claims cho JWT
	claims := jwt.MapClaims{
		"userid":   user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      expirationTime,
	}

	// Tạo token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Ký token với secretKey
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", 0, err
	}

	// Tính thời gian hết hạn
	expiresIn := int(expirationTime - time.Now().Unix())

	// Trả về token, thời gian hết hạn
	return tokenString, expiresIn, nil
}
