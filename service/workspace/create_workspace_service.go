package workspace

import (
	"api/dms"
	"api/notification"
	auth_utils "api/utils/auth"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/create_workspace_dtos"
	"github.com/timewise-team/timewise-models/models"
	"io/ioutil"
	"net/http"
)

type CreateWorkspaceService struct {
	// Thêm các dependencies cần thiết nếu có (ví dụ: database, API client, v.v.)
}

func NewCreateWorkspaceService() *CreateWorkspaceService {
	return &CreateWorkspaceService{}
}

// Main method
func (s *CreateWorkspaceService) InitWorkspace(workspaceRequest create_workspace_dtos.CreateWorkspaceRequest) (*create_workspace_dtos.CreateWorkspaceResponse, error) {
	//Validate workspace request
	if workspaceRequest.Email == "" {
		return nil, errors.New("Email is required")
	}
	if workspaceRequest.Title == "" {
		return nil, errors.New("Title is required")
	}
	if workspaceRequest.Description == "" {
		return nil, errors.New("Description is required")
	}
	if len(workspaceRequest.Title) > 50 {
		return nil, errors.New("Title must not exceed 50 characters")
	}
	if len(workspaceRequest.Description) > 255 {
		return nil, errors.New("Description must not exceed 255 characters")
	}
	if auth_utils.IsValidEmail(workspaceRequest.Email) == false {
		return nil, errors.New("Invalid email")
	}

	//Create workspace
	userEmail, err := s.GetUserEmailByEmail(workspaceRequest.Email)
	if err != nil {
		return nil, err
	}
	if userEmail == nil {
		return nil, err
	}
	var workspace = models.TwWorkspace{
		Title:       workspaceRequest.Title,
		Description: workspaceRequest.Description,
		Type:        "workspace",
		IsDeleted:   false,
	}
	workspaceResult, err := s.CreateWorkspace(workspace)
	if err != nil {
		return nil, err
	}
	if workspaceResult == nil {
		return nil, err
	}
	var workspaceUser = models.TwWorkspaceUser{
		WorkspaceId: workspaceResult.ID,
		Role:        "owner",
		Status:      "joined",
		IsActive:    true,
		IsVerified:  true,
		UserEmailId: userEmail.ID,
	}
	_, err = s.CreateWorkspaceUser(workspaceUser)
	if err != nil {
		return nil, err
	}

	// send notification
	notificationDto := models.TwNotifications{
		Title:       "New Workspace " + workspaceResult.Title + " created",
		Description: "You have created new workspace " + workspaceResult.Title + " successfully",
		Link:        fmt.Sprintf("/organization/%d", workspaceResult.ID),
		UserEmailId: userEmail.ID,
		Type:        "workspace_created",
		Message:     "You have created new workspace " + workspaceResult.Title + " successfully",
		IsSent:      true,
	}
	err = notification.PushNotifications(notificationDto)
	if err != nil {
		return nil, err
	}

	return workspaceResult, nil
}

func (s *CreateWorkspaceService) CreateWorkspaceUser(workspaceUser models.TwWorkspaceUser) (*models.TwWorkspaceUser, error) {
	//Validate workspace user
	if workspaceUser.WorkspaceId == 0 {
		return nil, errors.New("Workspace ID is required")
	}
	if workspaceUser.UserEmailId == 0 {
		return nil, errors.New("User Email ID is required")
	}
	if workspaceUser.Role == "" {
		return nil, errors.New("Role is required")
	}
	if workspaceUser.Status == "" {
		return nil, errors.New("Status is required")
	}
	if workspaceUser.Role != "owner" && workspaceUser.Role != "member" && workspaceUser.Role != "guest" && workspaceUser.Role != "admin" {
		return nil, errors.New("Invalid role")
	}
	if workspaceUser.Status != "joined" && workspaceUser.Status != "pending" {
		return nil, errors.New("Invalid status")
	}

	resp, err := dms.CallAPI(
		"POST",
		"/workspace_user",
		workspaceUser,
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

	var CreateWorkspaceUserResponse models.TwWorkspaceUser
	err = json.Unmarshal(body, &CreateWorkspaceUserResponse)
	if err != nil {
		return nil, err
	}
	return &CreateWorkspaceUserResponse, nil
}
func (s *CreateWorkspaceService) GetUserEmailByEmail(email string) (*models.TwUserEmail, error) {
	//Validate email
	if email == "" {
		return nil, errors.New("Email is required")
	}
	if auth_utils.IsValidEmail(email) == false {
		return nil, errors.New("Invalid email")
	}
	//Call API
	resp, err := dms.CallAPI(
		"GET",
		"/user_email/email/"+email,
		nil,
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

	var GetUserEmailByEmail models.TwUserEmail
	err = json.Unmarshal(body, &GetUserEmailByEmail)
	if err != nil {
		return nil, err
	}
	return &GetUserEmailByEmail, nil
}
func (s *CreateWorkspaceService) CreateWorkspace(workspace models.TwWorkspace) (*create_workspace_dtos.CreateWorkspaceResponse, error) {
	//Validate workspace
	if workspace.Title == "" {
		return nil, errors.New("Title is required")
	}
	if workspace.Description == "" {
		return nil, errors.New("Description is required")
	}
	if len(workspace.Title) > 50 {
		return nil, errors.New("Title must not exceed 50 characters")
	}
	if len(workspace.Description) > 255 {
		return nil, errors.New("Description must not exceed 255 characters")
	}
	if workspace.Type != "workspace" && workspace.Type != "personal" {
		return nil, errors.New("Invalid type")
	}
	if workspace.IsDeleted == true {
		return nil, errors.New("IsDeleted is required")
	}
	//Call API
	resp, err := dms.CallAPI(
		"POST",
		"/workspace",
		workspace,
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

	var CreateWorkspaceResponse create_workspace_dtos.CreateWorkspaceResponse
	err = json.Unmarshal(body, &CreateWorkspaceResponse)
	if err != nil {
		return nil, err
	}
	return &CreateWorkspaceResponse, nil
}
