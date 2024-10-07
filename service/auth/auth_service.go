package auth

import (
	"api/dms"
	"encoding/json"
	"github.com/timewise-team/timewise-models/models"
	"io/ioutil"
	"net/http"
)

type AuthService struct {
	// Thêm các dependencies cần thiết nếu có (ví dụ: database, API client, v.v.)
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) CreateUserEmail(user models.TwUser) (*models.TwUserEmail, error) {
	var UserEmail models.TwUserEmail
	UserEmail.Email = user.Email
	UserEmail.UserId = user.ID
	UserEmail.User = user

	resp, err := dms.CallAPI(
		"POST",
		"/user_email",
		UserEmail,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}

	var UserEmailResponse models.TwUserEmail
	err = json.Unmarshal(body, &UserEmailResponse)
	if err != nil {
		return nil, err
	}

	return &UserEmailResponse, nil
}

// CreateWorkspace handles creating the workspace via API
func (s *AuthService) CreateWorkspace() (*models.TwWorkspace, error) {
	var WorkspaceRequest models.TwWorkspace
	WorkspaceRequest.Title = "personal"
	WorkspaceRequest.Type = "personal"

	resp, err := dms.CallAPI(
		"POST",
		"/workspace",
		WorkspaceRequest,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}

	var WorkspaceResponse models.TwWorkspace
	err = json.Unmarshal(body, &WorkspaceResponse)
	if err != nil {
		return nil, err
	}

	return &WorkspaceResponse, nil
}
func (s *AuthService) CreateWorkspaceUser(userEmail *models.TwUserEmail, workspace *models.TwWorkspace) (*models.TwWorkspaceUser, error) {
	var WorkspaceUserRequest models.TwWorkspaceUser
	WorkspaceUserRequest.UserEmailId = userEmail.ID
	WorkspaceUserRequest.WorkspaceId = workspace.ID
	WorkspaceUserRequest.Workspace = *workspace
	WorkspaceUserRequest.UserEmail = *userEmail
	WorkspaceUserRequest.Role = "owner"
	WorkspaceUserRequest.Status = "joined"
	WorkspaceUserRequest.IsActive = true
	WorkspaceUserRequest.IsVerified = true
	WorkspaceUserRequest.ExtraData = ""
	WorkspaceUserRequest.WorkspaceKey = ""

	resp, err := dms.CallAPI(
		"POST",
		"/workspace_user",
		WorkspaceUserRequest,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}

	var WorkspaceUserResponse models.TwWorkspaceUser
	err = json.Unmarshal(body, &WorkspaceUserResponse)
	if err != nil {
		return nil, err
	}

	return &WorkspaceUserResponse, nil

}

func (s *AuthService) InitNewUser(user models.TwUser) (bool, error) {
	//_ = s.CreateNotificationSetting(user)
	// Create user email
	userEmailResponse, err := s.CreateUserEmail(user)
	if err != nil {
		return false, err // Return error if email creation fails
	}

	// Create workspace
	workspaceResponse, err := s.CreateWorkspace()
	if err != nil {
		return false, err // Return error if workspace creation fails
	}

	_, err = s.CreateWorkspaceUser(userEmailResponse, workspaceResponse)
	// Create workspace user

	return true, nil // Success
}

//func (s *AuthService) CreateNotificationSetting(user models.TwUser) *models.TwNotificationSettings {
//	var NotificationSetting models.TwNotificationSettings
//	NotificationSetting.UserId = user.ID
//	NotificationSetting.NotificationOnComment = true
//	NotificationSetting.NotificationOnDueDate = true
//	NotificationSetting.NotificationOnScheduleChange = true
//	NotificationSetting.NotificationOnDueDate = true
//
//	resp, err := dms.CallAPI(
//		"POST",
//		"/notification_setting",
//		NotificationSetting,
//		nil,
//		nil,
//		120,
//	)
//	if err != nil {
//		return nil
//	}
//
//	defer resp.Body.Close()
//
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil || resp.StatusCode != http.StatusOK {
//		return nil
//	}
//
//	var NotificationSettingResponse models.TwNotificationSettings
//	err = json.Unmarshal(body, &NotificationSettingResponse)
//	if err != nil {
//		return nil
//	}
//
//	return &NotificationSettingResponse
//
//}
