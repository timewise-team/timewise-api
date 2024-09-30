package auth_utils

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/timewise-team/timewise-models/models"
	"time"
)

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
