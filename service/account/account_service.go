package account

import (
	"api/dms"
	"encoding/json"
	"errors"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
	"github.com/timewise-team/timewise-models/models"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type AccountService struct {
}

func NewAccountService() *AccountService {
	return &AccountService{}
}

func parseResponseBody(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("unexpected status code: " + strconv.Itoa(resp.StatusCode))
	}
	return json.Unmarshal(body, v)
}

func (s *AccountService) GetUserInfoByUserId(userId string, status string) (core_dtos.GetUserResponseDto, error) {
	var userResponse models.TwUser
	resp, err := dms.CallAPI("GET", "/user/"+userId, nil, nil, nil, 120)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	if err := parseResponseBody(resp, &userResponse); err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}

	if userResponse.DeletedAt == nil {
		userResponse.DeletedAt = new(time.Time)
	}

	userDto := core_dtos.GetUserResponseDto{
		ID:                   userResponse.ID,
		CreatedAt:            userResponse.CreatedAt,
		UpdatedAt:            userResponse.UpdatedAt,
		DeteledAt:            *userResponse.DeletedAt,
		FirstName:            userResponse.FirstName,
		LastName:             userResponse.LastName,
		ProfilePicture:       userResponse.ProfilePicture,
		Timezone:             userResponse.Timezone,
		Locale:               userResponse.Locale,
		GoogleId:             userResponse.GoogleId,
		IsVerified:           userResponse.IsVerified,
		IsActive:             userResponse.IsActive,
		LastLoginAt:          userResponse.LastLoginAt,
		NotificationSettings: userResponse.NotificationSettings,
		CalendarSettings:     userResponse.CalendarSettings,
		Role:                 userResponse.Role,
	}

	emails, err := s.GetLinkedUserEmails(userId, status)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	userDto.Email = emails

	return userDto, nil
}

func (s *AccountService) GetLinkedUserEmails(userId string, status string) ([]core_dtos.EmailDto, error) {
	query := map[string]string{"status": status}
	resp, err := dms.CallAPI("GET", "/user_email/user/"+userId, nil, nil, query, 120)
	if err != nil {
		return nil, err
	}

	var userEmailResp []models.TwUserEmail
	if err := parseResponseBody(resp, &userEmailResp); err != nil {
		return nil, err
	}

	emails := make([]core_dtos.EmailDto, 0, len(userEmailResp))
	for _, email := range userEmailResp {
		emailStatus := ""
		if email.Status != nil {
			emailStatus = *email.Status
		}
		emails = append(emails, core_dtos.EmailDto{
			Email:  email.Email,
			Status: emailStatus,
		})
	}

	return emails, nil
}

func (s *AccountService) UpdateUserInfo(userId string, request core_dtos.UpdateProfileRequestDto) (core_dtos.GetUserResponseDto, error) {
	// call dms to update user info
	user := core_dtos.UpdateUserRequest{
		FirstName:            &request.FirstName,
		LastName:             &request.LastName,
		ProfilePicture:       &request.ProfilePicture,
		NotificationSettings: &request.NotificationSettings,
		CalendarSettings:     &request.CalendarSettings,
	}
	resp, err := dms.CallAPI("PUT", "/user/"+userId, user, nil, nil, 120)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}

	var userResp models.TwUser
	if err := parseResponseBody(resp, &userResp); err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}

	if userResp.DeletedAt == nil {
		userResp.DeletedAt = new(time.Time)
	}

	userDto := core_dtos.GetUserResponseDto{
		ID:                   userResp.ID,
		CreatedAt:            userResp.CreatedAt,
		UpdatedAt:            userResp.UpdatedAt,
		FirstName:            userResp.FirstName,
		LastName:             userResp.LastName,
		ProfilePicture:       userResp.ProfilePicture,
		Timezone:             userResp.Timezone,
		Locale:               userResp.Locale,
		GoogleId:             userResp.GoogleId,
		IsVerified:           userResp.IsVerified,
		IsActive:             userResp.IsActive,
		LastLoginAt:          userResp.LastLoginAt,
		NotificationSettings: userResp.NotificationSettings,
		CalendarSettings:     userResp.CalendarSettings,
		Role:                 userResp.Role,
	}

	emails, err := s.GetLinkedUserEmails(userId, "")
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	userDto.Email = emails

	return userDto, nil
}

func (s *AccountService) SendLinkAnEmailRequest(userId string, email string) (models.TwUserEmail, error) {
	// Check if the email is already a user.
	resp, err := dms.CallAPI("GET", "/user/get", nil, nil, map[string]string{"email": email}, 120)
	if err != nil {
		return models.TwUserEmail{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return models.TwUserEmail{}, errors.New("email is not a user")
	}

	// Check if email is already linked.
	resp, err = dms.CallAPI("GET", "/user_email/check", nil, nil, map[string]string{"email": email}, 120)
	if err != nil {
		return models.TwUserEmail{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		return models.TwUserEmail{}, errors.New("email is already linked or pending")
	}

	// Fetch user email details.
	resp, err = dms.CallAPI("GET", "/user_email/email/"+email, nil, nil, nil, 120)
	if err != nil {
		return models.TwUserEmail{}, err
	}
	defer resp.Body.Close()

	var userEmail models.TwUserEmail
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return models.TwUserEmail{}, errors.New("failed to fetch email details")
	}
	if err := json.Unmarshal(body, &userEmail); err != nil {
		return models.TwUserEmail{}, err
	}

	// Update status to "pending."
	query := map[string]string{"email": email, "status": "pending", "target_user_id": userId}
	if _, err := dms.CallAPI("PATCH", "/user_email/status", nil, nil, query, 120); err != nil {
		return models.TwUserEmail{}, err
	}

	return userEmail, nil
}

func (s *AccountService) UpdateStatusLinkEmailRequest(userId string, email string, status string) (core_dtos.GetUserResponseDto, error) {
	queryParams := map[string]string{
		"email":  email,
		"status": status,
	}
	respEmail, err := dms.CallAPI("PATCH", "/user_email", nil, nil, queryParams, 120)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	defer respEmail.Body.Close()
	if respEmail.StatusCode != http.StatusOK {
		return core_dtos.GetUserResponseDto{}, errors.New("cannot update status of email")
	}
	// return user info
	return s.GetUserInfoByUserId(userId, "")
}

func (s *AccountService) UnlinkAnEmail(email string) (core_dtos.GetUserResponseDto, error) {
	// check if email is already is linked to a user
	respUserEmail, err := dms.CallAPI("GET", "/user_email/email/"+email, nil, nil, nil, 120)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	defer respUserEmail.Body.Close()
	body, err := ioutil.ReadAll(respUserEmail.Body)
	if err != nil || respUserEmail.StatusCode != http.StatusOK {
		return core_dtos.GetUserResponseDto{}, err
	}
	// marshal response body
	var userEmailResp models.TwUserEmail
	err = json.Unmarshal(body, &userEmailResp)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	if userEmailResp.Status == nil || *userEmailResp.Status != "linked" {
		return core_dtos.GetUserResponseDto{}, errors.New("Email is not linked to any user")
	}
	// call dms to get user_id by email in user_email
	queryParam := map[string]string{
		"email": email,
	}
	resp, err := dms.CallAPI("GET", "/user/get", nil, nil, queryParam, 120)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return core_dtos.GetUserResponseDto{}, errors.New("Can not get user_id by email")
	}
	// marshal response body
	var usersResp models.TwUser
	err = json.Unmarshal(body, &usersResp)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	userId := usersResp.ID
	userIdStr := strconv.Itoa(userId)
	// call dms to change current user_id to user_id got from above api in user_email
	queryParams := map[string]string{
		"email":  email,
		"status": "",
	}
	_, err = dms.CallAPI("PATCH", "/user_email", nil, nil, queryParams, 120)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	// return user info
	return s.GetUserInfoByUserId(userIdStr, "")
}

func (s *AccountService) DeactivateAccount(userId string) error {
	isActive := false
	request := core_dtos.UpdateUserRequest{
		IsActive: &isActive,
	}
	// call dms to deactivate account
	_, err := dms.CallAPI("PUT", "/user/"+userId, request, nil, nil, 120)
	if err != nil {
		return err
	}
	return nil
}
