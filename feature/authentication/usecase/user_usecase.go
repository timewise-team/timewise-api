package usecase

import (
	"api/utils/auth"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/user_login_dtos"
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
func CallAPI(method string, url string, body interface{}, headers map[string]string, queryParams map[string]string, timeout time.Duration) (*http.Response, error) {
	// Chuyển body thành JSON nếu body không rỗng
	var requestBody []byte
	var err error
	if body != nil {
		requestBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	// Tạo HTTP client với timeout
	client := &http.Client{
		Timeout: timeout,
	}

	// Tạo request với method (GET, POST, PUT, etc.)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	// Thêm query parameters nếu có
	if len(queryParams) > 0 {
		q := req.URL.Query()
		for key, value := range queryParams {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	// Thiết lập headers cho request nếu có
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Gửi request và nhận phản hồi
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Trả về response để hàm gọi xử lý tiếp
	return resp, nil
}

func Login(req user_login_dtos.UserLoginRequest) (*user_login_dtos.UserLoginResponse, error) {

	// Sử dụng hàm CallAPI để gọi API DMS
	resp, err := CallAPI("POST", "user/login", req, nil, nil, 10*time.Second)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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

	var user user_login_dtos.UserLoginRequest
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
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
