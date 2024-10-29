package account

import (
	"api/dms"
	"encoding/json"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
	"github.com/timewise-team/timewise-models/models"
	"io/ioutil"
	"net/http"
)

type AccountService struct {
}

func NewAccountService() *AccountService {
	return &AccountService{}
}

func (h *AccountService) GetUserInfoByUserId(userId string) (core_dtos.GetUserResponseDto, error) {
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

func (h *AccountService) UpdateUserInfo(userId string, request core_dtos.UpdateProfileRequestDto) (core_dtos.GetUserResponseDto, error) {
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

func (h *AccountService) GetLinkedUserEmails(userId string) ([]string, error) {
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
