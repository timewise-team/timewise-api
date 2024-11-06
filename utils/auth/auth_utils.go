package auth_utils

import (
	"api/config"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/timewise-team/timewise-models/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

type GoogleOauthData struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
}

var cfg, err = config.LoadConfig()

var GoogleOauth = oauth2.Config{
	ClientID:     cfg.GoogleOauth.ClientID,
	ClientSecret: cfg.GoogleOauth.ClientSecret,
	RedirectURL:  "",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

type TokenInfo struct {
	ExpiresIn int64 `json:"expires_in"`
	// other fields as necessary
}

func VerifyGoogleToken(code string) ([]byte, error) {
	response, err := http.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token=" + code)
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
func CheckGoogleTokenExpiry(accessToken string) error {
	response, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v1/tokeninfo?access_token=%s", accessToken))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New("invalid or expired token")
	}

	var tokenInfo TokenInfo
	if err := json.NewDecoder(response.Body).Decode(&tokenInfo); err != nil {
		return err
	}

	// Check if the token is expired
	if tokenInfo.ExpiresIn <= 0 {
		return errors.New("token has expired")
	}

	return nil
}
func GenerateJWTToken(user models.TwUser, secretKey string) (string, int, error) {
	// Định nghĩa thời gian hết hạn cho token (ví dụ: 2 giờ)
	expirationTime := time.Now().Add(168 * time.Hour).Unix()

	// Tạo claims cho JWT
	claims := jwt.MapClaims{
		"userid": user.ID,
		"email":  user.Email,
		"exp":    expirationTime,
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

func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func GenerateInvitationToken(workspaceId int, action string, secretKey string, email string, role string) (string, error) {
	claims := jwt.MapClaims{
		"email":        email,
		"workspace_id": workspaceId,
		"role":         role,
		"action":       action,                                // accept hoặc decline
		"exp":          time.Now().Add(24 * time.Hour).Unix(), // Token có thời hạn 24h
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateScheduleInvitationToken(workspaceUserId int, action string, secretKey string, scheduleId int) (string, error) {
	claims := jwt.MapClaims{
		"schedule_id":       scheduleId,
		"workspace_user_id": workspaceUserId,
		"action":            action,                                // accept hoặc decline
		"exp":               time.Now().Add(24 * time.Hour).Unix(), // Token có thời hạn 24h
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func compressData(data []byte) (string, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write(data)
	if err != nil {
		return "", err
	}
	err = gz.Close()
	if err != nil {
		return "", err
	}

	// Base64 encode the compressed data for safe transmission
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
func decompressData(compressedData string) (string, error) {
	// Decode dữ liệu từ base64
	data, err := base64.StdEncoding.DecodeString(compressedData)
	if err != nil {
		return "", err
	}

	// Sử dụng gzip để giải nén
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer reader.Close()

	// Đọc và trả về dữ liệu đã giải nén
	decompressedData, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(decompressedData), nil
}
func ParseInvitationToken(tokenString string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(jwt.MapClaims), nil
}
func ParseLinkEmailToken(tokenString string, secretKey string) (jwt.MapClaims, error) {
	// Giải mã token với key bí mật
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	// Lấy claims của token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	// Lấy dữ liệu đã được nén trong token (trong trường hợp của chúng ta là 'data')
	compressedData, ok := claims["data"].(string)
	if !ok {
		return nil, errors.New("missing or invalid data in token")
	}

	// Giải nén dữ liệu
	decompressedData, err := decompressData(compressedData)
	if err != nil {
		return nil, err
	}

	// Sau khi giải nén, bạn cần chuyển lại thành MapClaims (hoặc dữ liệu cần thiết)
	var decompressedClaims map[string]interface{}
	err = json.Unmarshal([]byte(decompressedData), &decompressedClaims)
	if err != nil {
		return nil, err
	}

	// Chuyển các dữ liệu đã giải nén thành jwt.MapClaims và trả về
	return decompressedClaims, nil
}
func GenerateLinkEmailToken(currentUid string, email string, action string, secretKey string) (string, error) {
	// Prepare payload data
	payload := fmt.Sprintf("{\"uid\":\"%s\",\"email\":\"%s\",\"action\":\"%s\",\"exp\":%d}", currentUid, email, action, time.Now().Add(24*time.Hour).Unix())

	// Compress the payload
	compressedPayload, err := compressData([]byte(payload))
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"data": compressedPayload,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
