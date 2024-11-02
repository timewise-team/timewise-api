package workspace

import (
	"api/dms"
	"encoding/json"
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
	return workspaceResult, nil
}

func (s *CreateWorkspaceService) CreateWorkspaceUser(workspaceUser models.TwWorkspaceUser) (*models.TwWorkspaceUser, error) {

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
