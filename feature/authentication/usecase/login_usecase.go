package usecase

import (
	"api/dms"
	"api/utils/auth"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/user_login_dtos"
	"github.com/timewise-team/timewise-models/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func Login(req user_login_dtos.UserLoginRequest) (*user_login_dtos.UserLoginResponse, error) {
	// Sử dụng hàm CallAPI để gọi API DMS
	resp, err := dms.CallAPI("POST", "/user/login", req, nil, nil, 10*time.Second)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Printf("resp: %v\n", resp.StatusCode)
	// Xử lý phản hồi từ API
	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, errors.New("failed to decode error response")
		}
		if message, ok := errorResponse["message"].(string); ok {
			return nil, errors.New(message)
		}
		return nil, errors.New("unknown error occurred")
	}

	var user models.TwUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("Invalid email or password")
	}
	accessToken, expiresIn, err := auth_utils.GenerateJWTToken(user, viper.GetString("JWT_SECRET"))
	if err != nil {
		return nil, err
	}

	// Tạo đối tượng UserLoginResponse
	response := &user_login_dtos.UserLoginResponse{
		AccessToken: accessToken,
		ExpiresIn:   expiresIn,
		TokenType:   "Bearer", // Sử dụng "Bearer" cho token loại này

	}

	return response, nil
}
