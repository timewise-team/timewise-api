package account

import (
	"api/dms"
	auth_service "api/service/auth"
	auth_utils "api/utils/auth"
	"encoding/json"
	"errors"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
	dtos "github.com/timewise-team/timewise-models/dtos/core_dtos/user_register_dtos"
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

func (s *AccountService) GetLinkedUserEmails(userId string) ([]string, error) {
	resp, err := dms.CallAPI("GET", "/user_email/user/"+userId, nil, nil, nil, 120)
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

func (s *AccountService) LinkAnEmail(userId string, oauthData auth_utils.GoogleOauthData) (core_dtos.GetUserResponseDto, error) {
	getOrCreateUserReq := dtos.GetOrCreateUserRequestDto{
		Email:          oauthData.Email,
		FullName:       oauthData.Name,
		ProfilePicture: oauthData.Picture,
		VerifiedEmail:  oauthData.VerifiedEmail,
		GoogleId:       oauthData.Id,
		GivenName:      oauthData.GivenName,
		FamilyName:     oauthData.FamilyName,
		Locale:         oauthData.Locale,
	}

	// Get or Create user
	resp, err := dms.CallAPI(
		"POST",
		"/user/get-create",
		getOrCreateUserReq,
		nil,
		nil,
		120,
	)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, errors.New("could not get or create user")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return core_dtos.GetUserResponseDto{}, err
	}

	// marshal response body
	var userRespDto dtos.GetOrCreateUserResponseDto
	err = json.Unmarshal(body, &userRespDto)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}

	if userRespDto.IsNewUser {
		_, err := auth_service.NewAuthService().InitNewUser(userRespDto.User)
		if err != nil {
			return core_dtos.GetUserResponseDto{}, errors.New("could not init new user")
		}
	}
	queryParams := map[string]string{
		"user_id": userId,
		"email":   oauthData.Email,
	}
	// else then update user_id to user_email
	respEmail, err := dms.CallAPI("PATCH", "/user_email", nil, nil, queryParams, 120)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(respEmail.Body)
	if err != nil || respEmail.StatusCode != http.StatusOK {
		return core_dtos.GetUserResponseDto{}, err
	}

	// marshal response body
	var userEmailResp models.TwUserEmail
	err = json.Unmarshal(body, &userEmailResp)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}

	// return user info
	resp, err = dms.CallAPI("GET", "/user/"+userId, nil, nil, nil, 120)
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
	userEmailList, err := s.GetLinkedUserEmails(userId)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	userDto.Email = userEmailList
	return userDto, nil
}

func (s *AccountService) UnlinkAnEmail(email string) (core_dtos.GetUserResponseDto, error) {
	// call dms to get user_id by email in user_email
	queryParam := map[string]string{
		"email": email,
	}
	resp, err := dms.CallAPI("GET", "/user", nil, nil, queryParam, 120)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return core_dtos.GetUserResponseDto{}, err
	}
	// marshal response body
	var userEmailResp []models.TwUser
	err = json.Unmarshal(body, &userEmailResp)
	if err != nil {
		return core_dtos.GetUserResponseDto{}, err
	}
	userId := userEmailResp[0].ID
	userIdStr := strconv.Itoa(userId)
	// call dms to change current user_id to user_id got from above api in user_email
	queryParams := map[string]string{
		"user_id": userIdStr,
		"email":   email,
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
	userEmailList, err := s.GetLinkedUserEmails(userIdStr)
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
