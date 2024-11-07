package account

import (
	"api/dms"
	"encoding/json"
	"errors"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
	"github.com/timewise-team/timewise-models/models"
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

func (s *AccountService) GetUserInfoByUserId(userId string) (core_dtos.GetUserResponseDto, error) {
	// call dms to query database
	resp, err := dms.CallAPI("GET", "/user/"+userId, nil, nil, nil, 120)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return core_dtos.GetUserResponseDto{}, err
	}
	// marshal response body
	var userResponse models.TwUser
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	// parse userResponse to userDto
	if userResponse.DeletedAt == nil {
		userResponse.DeletedAt = new(time.Time)
	}
	userDto := core_dtos.GetUserResponseDto{
		ID:        userResponse.ID,
		CreatedAt: userResponse.CreatedAt,
		UpdatedAt: userResponse.UpdatedAt,
		DeteledAt: *userResponse.DeletedAt,
		FirstName: userResponse.FirstName,
		LastName:  userResponse.LastName,
		//Email:  (need to call another API to get all email of this user)
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

	resp, err = dms.CallAPI("GET", "/user_email/user/"+userId, nil, nil, nil, 120)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return core_dtos.GetUserResponseDto{}, err
	}
	// marshal response body
	var userEmailResp []models.TwUserEmail
	err = json.Unmarshal(body, &userEmailResp)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	// parse userEmailResp to []string
	emailSlice := make([]string, 0)
	for _, email := range userEmailResp {
		emailSlice = append(emailSlice, email.Email)
	}
	userDto.Email = emailSlice
	// return user info
	return userDto, nil
}

func (s *AccountService) UpdateUserInfo(userId string, request core_dtos.UpdateProfileRequestDto) (core_dtos.GetUserResponseDto, error) {
	// call dms to update user info
	user := models.TwUser{
		FirstName:            request.FirstName,
		LastName:             request.LastName,
		ProfilePicture:       request.ProfilePicture,
		NotificationSettings: request.NotificationSettings,
		CalendarSettings:     request.CalendarSettings,
	}
	resp, err := dms.CallAPI("PUT", "/user/"+userId, user, nil, nil, 120)
	defer resp.Body.Close()
	var userResp models.TwUser
	// unmarshal response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return core_dtos.GetUserResponseDto{}, err
	}
	err = json.Unmarshal(body, &userResp)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	// parse userResp to userDto
	userDto := core_dtos.GetUserResponseDto{
		ID:        userResp.ID,
		CreatedAt: userResp.CreatedAt,
		UpdatedAt: userResp.UpdatedAt,
		FirstName: userResp.FirstName,
		LastName:  userResp.LastName,
		//Email:  (need to call another API to get all email of this user)
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
	resp, err = dms.CallAPI("GET", "/user_email/user/"+userId, nil, nil, nil, 120)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return core_dtos.GetUserResponseDto{}, err
	}
	// marshal response body
	var userEmailResp []models.TwUserEmail
	err = json.Unmarshal(body, &userEmailResp)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	// parse userEmailResp to []string
	emailSlice := make([]string, 0)
	for _, email := range userEmailResp {
		emailSlice = append(emailSlice, email.Email)
	}
	userDto.Email = emailSlice
	return userDto, nil
}

func (s *AccountService) GetLinkedUserEmails(userId string, status string) ([]string, error) {
	query := map[string]string{
		"status": status,
	}
	resp, err := dms.CallAPI("GET", "/user_email/user/"+userId, nil, nil, query, 120)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}
	// marshal response body
	var userEmailResp []models.TwUserEmail
	err = json.Unmarshal(body, &userEmailResp)
	if err != nil {
		return nil, err
	}
	// parse userEmailResp to []string
	emailSlice := make([]string, 0)
	for _, email := range userEmailResp {
		emailSlice = append(emailSlice, email.Email)
	}
	return emailSlice, nil
}

func (s *AccountService) SendLinkAnEmailRequest(userId string, email string) (models.TwUserEmail, error) {
	// check if email is already is a user
	resp, err := dms.CallAPI("GET", "/user/get", nil, nil, map[string]string{"email": email}, 120)
	if err != nil {
		return models.TwUserEmail{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return models.TwUserEmail{}, errors.New("Email is not already a user. Only existing user can be linked")
	}
	if err != nil || resp.StatusCode != http.StatusOK {
		return models.TwUserEmail{}, errors.New("Cannot check if email is already a user")
	}
	// get user info from user_emails table
	resp, err = dms.CallAPI("GET", "/user_email/check", nil, nil, map[string]string{"email": email, "user_id": userId}, 120)
	if err != nil {
		return models.TwUserEmail{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		return models.TwUserEmail{}, errors.New("Email is already linked or rejected or pending")
	}
	resp, err = dms.CallAPI("GET", "/user_email/email/"+email, nil, nil, nil, 120)
	if err != nil {
		return models.TwUserEmail{}, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return models.TwUserEmail{}, errors.New("Cannot check fetch user email info")
	}
	// marshal response body
	var userEmailResp models.TwUserEmail
	err = json.Unmarshal(body, &userEmailResp)
	if err != nil {
		return models.TwUserEmail{}, err
	}
	//if userEmailResp.Status != nil {
	//	return models.TwUserEmail{}, errors.New("Email is already linked or rejected or pending")
	//}
	// call dms to create a new user_email
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return models.TwUserEmail{}, err
	}
	status := "pending"
	userEmail := models.TwUserEmail{
		UserId: userIdInt,
		Email:  email,
		Status: &status,
	}
	resp, err = dms.CallAPI("POST", "/user_email", userEmail, nil, nil, 120)
	if err != nil {
		return models.TwUserEmail{}, err
	}

	return userEmailResp, nil
}

func (s *AccountService) UpdateStatusLinkEmailRequest(userId string, email string, status string) (core_dtos.GetUserResponseDto, error) {
	queryParams := map[string]string{
		"user_id": userId,
		"email":   email,
		"status":  status,
	}
	// delete pending email if status is rejected or accepted
	if status == "linked" {
		queryParams := map[string]string{
			"user_id": userId,
			"email":   email,
			"status":  "pending",
		}
		respEmail, err := dms.CallAPI("DELETE", "/user_email", nil, nil, queryParams, 120)
		if err != nil {
			return core_dtos.GetUserResponseDto{}, err
		}
		defer respEmail.Body.Close()
		if respEmail.StatusCode != http.StatusOK {
			return core_dtos.GetUserResponseDto{}, errors.New("cannot delete email")
		}
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
	resp, err := dms.CallAPI("GET", "/user/"+userId, nil, nil, nil, 120)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return core_dtos.GetUserResponseDto{}, err
	}
	// marshal response body
	var userResponse models.TwUser
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	if userResponse.DeletedAt == nil {
		userResponse.DeletedAt = new(time.Time)
	}
	// parse userResponse to userDto
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
	userEmailList, err := s.GetLinkedUserEmails(userId, "")
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	userDto.Email = userEmailList
	return userDto, nil
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
		"user_id": userIdStr,
		"email":   email,
		"status":  "",
	}
	_, err = dms.CallAPI("PATCH", "/user_email", nil, nil, queryParams, 120)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	// return user info
	resp, err = dms.CallAPI("GET", "/user/"+userIdStr, nil, nil, nil, 120)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return core_dtos.GetUserResponseDto{}, err
	}
	// marshal response body
	var userResponse models.TwUser
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	if userResponse.DeletedAt == nil {
		userResponse.DeletedAt = new(time.Time)
	}
	// parse userResponse to userDto
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
	userEmailList, err := s.GetLinkedUserEmails(userIdStr, "")
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	userDto.Email = userEmailList
	return userDto, nil
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
