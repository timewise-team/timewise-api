package dms

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

type DmsHandler struct {
	DMSUrl string
}

func CallAPI(method string, path string, body interface{}, headers map[string]string, queryParams map[string]string, timeout time.Duration) (*http.Response, error) {
	dmsUrl := viper.GetString("DMS_URL")

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
	client := &http.Client{}

	// Tạo request với method (GET, POST, PUT, etc.)
	url := dmsUrl + path
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
