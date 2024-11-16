package unit_test_test

import (
	"api/service/account"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
	"github.com/timewise-team/timewise-models/models"
	"io/ioutil"
	"net/http"
	"testing"
)

// MockDMSClient to simulate the dms package's CallAPI function
type MockDMSClient struct {
	mock.Mock
}

func (m *MockDMSClient) CallAPI(method, url string, body interface{}, headers map[string]string, query map[string]string, timeout int) (*http.Response, error) {
	args := m.Called(method, url, body, headers, query, timeout)
	resp, _ := args.Get(0).(*http.Response)
	return resp, args.Error(1)
}

// Helper function to create a mocked HTTP response
func newMockResponse(statusCode int, body interface{}) *http.Response {
	jsonBody, _ := json.Marshal(body)
	return &http.Response{
		StatusCode: statusCode,
		Body:       ioutil.NopCloser(bytes.NewReader(jsonBody)),
	}
}

func TestGetUserInfoByUserId(t *testing.T) {
	mockDMS := new(MockDMSClient)
	service := account.NewAccountService()

	// Mock response for /user/{userId}
	mockUser := models.TwUser{
		ID:             1,
		FirstName:      "John",
		LastName:       "Doe",
		ProfilePicture: "url/to/pic",
		IsActive:       true,
	}
	mockDMS.On("CallAPI", "GET", "/user/1", nil, nil, nil, 120).
		Return(newMockResponse(http.StatusOK, mockUser), nil)

	// Mock response for /user_email/user/{userId}
	mockEmails := []models.TwUserEmail{
		{Email: "john.doe@example.com", Status: stringPtr("linked")},
	}
	mockDMS.On("CallAPI", "GET", "/user_email/user/1", nil, nil, map[string]string{"status": ""}, 120).
		Return(newMockResponse(http.StatusOK, mockEmails), nil)

	// Test the function
	userDto, err := service.GetUserInfoByUserId("1", "")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, mockUser.FirstName, userDto.FirstName)
	assert.Equal(t, 1, userDto.ID)
	assert.Len(t, userDto.Email, 1)
	assert.Equal(t, "john.doe@example.com", userDto.Email[0].Email)

	mockDMS.AssertExpectations(t)
}

func TestUpdateUserInfo(t *testing.T) {
	mockDMS := new(MockDMSClient)
	service := account.NewAccountService()

	// Mock request body
	request := core_dtos.UpdateProfileRequestDto{
		FirstName:      "Jane",
		LastName:       "Smith",
		ProfilePicture: "url/to/new-pic",
	}

	// Mock response for /user/{userId}
	updatedUser := models.TwUser{
		ID:             1,
		FirstName:      "Jane",
		LastName:       "Smith",
		ProfilePicture: "url/to/new-pic",
	}
	mockDMS.On("CallAPI", "PUT", "/user/1", mock.Anything, nil, nil, 120).
		Return(newMockResponse(http.StatusOK, updatedUser), nil)

	// Test the function
	userDto, err := service.UpdateUserInfo("1", request)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "Jane", userDto.FirstName)
	assert.Equal(t, "Smith", userDto.LastName)
	assert.Equal(t, "url/to/new-pic", userDto.ProfilePicture)

	mockDMS.AssertExpectations(t)
}

func TestDeactivateAccount(t *testing.T) {
	mockDMS := new(MockDMSClient)
	service := account.NewAccountService()

	// Mock response for /user/{userId}
	mockDMS.On("CallAPI", "PUT", "/user/1", mock.Anything, nil, nil, 120).
		Return(newMockResponse(http.StatusOK, nil), nil)

	// Test the function
	err := service.DeactivateAccount("1")

	// Assertions
	assert.NoError(t, err)

	mockDMS.AssertExpectations(t)
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
