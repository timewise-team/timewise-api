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

func (s AccountService) convertToUserResponseDto(userResponse models.TwUser, status string) core_dtos.GetUserResponseDto {
	if userResponse.DeletedAt == nil {
		userResponse.DeletedAt = new(time.Time)
	}
	emails, err := s.GetLinkedUserEmails(strconv.Itoa(userResponse.ID), status)
	if err != nil {
		return core_dtos.GetUserResponseDto{}
	}
	return core_dtos.GetUserResponseDto{
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
		Email:                emails,
	}
}

func (s *AccountService) updateEmailStatus(email, userId, status string) error {
	if status == "" {
		userId = ""
	}
	query := map[string]string{"email": email, "status": status, "target_user_id": userId}
	_, err := dms.CallAPI("PATCH", "/user_email/status", nil, nil, query, 120)
	return err
}

func (s *AccountService) checkIfUserExists(email string) (bool, error) {
	resp, err := dms.CallAPI("GET", "/user/get", nil, nil, map[string]string{"email": email}, 120)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return true, nil
}

func (s *AccountService) checkIfEmailLinked(email string) (bool, error) {
	resp, err := dms.CallAPI("GET", "/user_email/check", nil, nil, map[string]string{"email": email}, 120)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		return true, nil // Email đã liên kết
	}

	return false, nil
}

func (s *AccountService) getEmailDetails(email string) (models.TwUserEmail, error) {
	resp, err := dms.CallAPI("GET", "/user_email/email/"+email, nil, nil, nil, 120)
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

	return userEmail, nil
}

func (s *AccountService) fetchUserByEmail(email string) (models.TwUser, error) {
	resp, err := dms.CallAPI("GET", "/user/get", nil, nil, map[string]string{"email": email}, 120)
	if err != nil {
		return models.TwUser{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.TwUser{}, errors.New("failed to fetch user by email")
	}

	var user models.TwUser
	if err := parseResponseBody(resp, &user); err != nil {
		return models.TwUser{}, err
	}

	return user, nil
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

	return s.convertToUserResponseDto(userResponse, status), nil
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
			ID:     email.ID,
			Email:  email.Email,
			Status: emailStatus,
		})
	}

	return emails, nil
}

func (s *AccountService) UpdateUserInfo(userId string, request core_dtos.UpdateProfileRequestDto) (core_dtos.GetUserResponseDto, error) {
	requestBody := core_dtos.UpdateUserRequest{
		FirstName:      &request.FirstName,
		LastName:       &request.LastName,
		ProfilePicture: &request.ProfilePicture,
	}
	// call dms to update user info
	var userResponse models.TwUser
	resp, err := dms.CallAPI("PUT", "/user/"+userId, requestBody, nil, nil, 120)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}

	if err := parseResponseBody(resp, &userResponse); err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}

	return s.convertToUserResponseDto(userResponse, "pending"), nil
}

func (s *AccountService) SendLinkAnEmailRequest(userId string, targetEmail string, currentEmail string) (models.TwUserEmail, error) {
	// Check if the email is already a user.
	isUser, err := s.checkIfUserExists(targetEmail)
	if err != nil {
		return models.TwUserEmail{}, err
	}
	if !isUser {
		return models.TwUserEmail{}, errors.New("email is not a user")
	}
	// Check if email is already linked.
	isLinked, err := s.checkIfEmailLinked(targetEmail)
	if err != nil {
		return models.TwUserEmail{}, err
	}
	if isLinked {
		return models.TwUserEmail{}, errors.New("email is already linked or pending")
	}

	// check if email is request to link to parents
	parentEmail, err := s.GetParentLinkedEmails(currentEmail)
	if err != nil {
		return models.TwUserEmail{}, err
	}
	if parentEmail != "" && parentEmail == targetEmail {
		return models.TwUserEmail{}, errors.New("you are already be linked to: " + parentEmail)
	}

	// Update status to "pending."
	err = s.updateEmailStatus(targetEmail, userId, "pending")
	if err != nil {
		return models.TwUserEmail{}, err
	}
	// Fetch user email details.
	userEmail, err := s.getEmailDetails(targetEmail)
	if err != nil {
		return models.TwUserEmail{}, err
	}

	return userEmail, nil
}

func (s *AccountService) UpdateStatusLinkEmailRequest(userId string, email string, status string) (core_dtos.GetUserResponseDto, error) {
	if status != "pending" && status != "linked" && status != "rejected" && status != "" {
		return s.GetUserInfoByUserId(userId, "invalid status")
	}
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
	userEmail, err := s.getEmailDetails(email)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	if userEmail.Status == nil || *userEmail.Status != "linked" {
		return core_dtos.GetUserResponseDto{}, errors.New("email is not linked to any user")
	}
	// call dms to get user_id by email in user_email
	user, err := s.fetchUserByEmail(email)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	if err := s.updateEmailStatus(email, strconv.Itoa(user.ID), ""); err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	// return user info
	return s.GetUserInfoByUserId(strconv.Itoa(user.ID), "")
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

func (s *AccountService) GetUserByUserId(userId string) (models.TwUser, error) {
	resp, err := dms.CallAPI("GET", "/user/"+userId, nil, nil, nil, 120)
	if err != nil {
		return models.TwUser{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.TwUser{}, errors.New("failed to get user")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.TwUser{}, err
	}
	var user models.TwUser
	if err := json.Unmarshal(body, &user); err != nil {
		return models.TwUser{}, err
	}
	return user, nil

}

func (s *AccountService) GetParentLinkedEmails(email string) (string, error) {
	resp, err := dms.CallAPI("GET", "/user_email/email/"+email, nil, nil, nil, 120)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to get parent linked email")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var currentIsLinkTo models.TwUserEmail
	if err := json.Unmarshal(body, &currentIsLinkTo); err != nil {
		return "", err
	}
	if currentIsLinkTo.IsLinkedTo == nil {
		return "", nil
	}
	resp, err = dms.CallAPI("GET", "/user_email/user_id/"+strconv.Itoa(*currentIsLinkTo.IsLinkedTo), nil, nil, nil, 120)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to get parent linked email")
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var parentEmail models.TwUserEmail
	if err := json.Unmarshal(body, &parentEmail); err != nil {
		return "", err
	}
	return parentEmail.Email, nil
}
