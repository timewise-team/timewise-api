package usecase

import (
	"api/config"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/user_login_dtos"
	user_register_dto "github.com/timewise-team/timewise-models/dtos/core_dtos/user_register_dtos"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"strings"
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

func CallDMSAPIForUserLogin(req user_login_dtos.UserLoginRequest, cfg *config.Config) (*user_login_dtos.UserLoginResponse, error) {
	// Sử dụng hàm CallAPI để gọi API DMS
	resp, err := CallAPI("POST", cfg.BaseURL+"user/login", req, nil, nil, 10*time.Second)
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

	accessToken, expiresIn, err := GenerateJWTToken(user, cfg.JWT_SECRET)
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

func CallDMSAPIForRegister(RegisterRequestDto user_register_dto.RegisterRequestDto, cfg *config.Config) error {

	// Check if passwords match
	if RegisterRequestDto.Password != RegisterRequestDto.ConfirmPassword {
		return errors.New("Passwords do not match")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(RegisterRequestDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Error hashing password")
	}

	fullName := strings.TrimSpace(RegisterRequestDto.FullName)
	lastSpaceIndex := strings.LastIndex(fullName, " ")

	var firstName, lastName string
	if lastSpaceIndex != -1 {
		firstName = fullName[:lastSpaceIndex]
		lastName = fullName[lastSpaceIndex+1:]
	} else {
		firstName = fullName
		lastName = ""
	}

	registerResponse := user_register_dto.RegisterResponseDto{
		UserName:     RegisterRequestDto.UserName,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        RegisterRequestDto.Email,
		HashPassword: string(hashedPassword),
	}

	resp, err := CallAPI("POST", cfg.BaseURL+"auth/register", registerResponse, nil, nil, 10*time.Second)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Can not read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}

	return nil
}
