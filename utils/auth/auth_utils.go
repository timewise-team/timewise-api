package auth_utils

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/user_login_dtos"
	"time"
)

func GenerateJWTToken(user user_login_dtos.UserLoginRequest, secretKey string) (string, int, error) {
	// Định nghĩa thời gian hết hạn cho token (ví dụ: 2 giờ)
	expirationTime := time.Now().Add(2 * time.Hour).Unix()

	// Tạo claims cho JWT
	claims := jwt.MapClaims{
		"username": user.Username,
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
