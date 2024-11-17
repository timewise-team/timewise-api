package workspace

import (
	"api/dms"
	auth_utils "api/utils/auth"
	"encoding/json"
	"errors"
	"github.com/timewise-team/timewise-models/models"
	"io/ioutil"
	"net/http"
	"strconv"
)

type WorkspaceService struct {
	// Thêm các dependencies cần thiết nếu có (ví dụ: database, API client, v.v.)
}

func NewWorkspaceService() *WorkspaceService {
	return &WorkspaceService{}
}

func (s *WorkspaceService) GetWorkspacesByEmail(email string) ([]models.TwWorkspace, error) {
	// Validate email
	if email == "" {
		return nil, errors.New("Invalid email")
	}
	if !auth_utils.IsValidEmail(email) {
		return nil, errors.New("Invalid email")
	}
	// Call API
	resp, err := dms.CallAPI(
		"GET",
		"/workspace/email/"+email,
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

	var workspaces []models.TwWorkspace
	err = json.Unmarshal(body, &workspaces)
	if err != nil {
		return nil, err
	}

	return workspaces, nil
}

func (s *WorkspaceService) GetWorkspacesByUserId(userId string) ([]models.TwWorkspace, error) {
	_, err := strconv.Atoi(userId)
	if err != nil {
		return nil, errors.New("Invalid user id")
	}
	// Call API
	resp, err := dms.CallAPI(
		"GET",
		"/workspace/user/"+userId,
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

	var workspaces []models.TwWorkspace
	err = json.Unmarshal(body, &workspaces)
	if err != nil {
		return nil, err
	}

	return workspaces, nil
}

func (s *WorkspaceService) GetWorkspaceById(workspaceId string) *models.TwWorkspace {
	// Validate workspaceId
	if workspaceId == "" {
		return nil
	}
	//Check is number workspaceId
	_, err := strconv.Atoi(workspaceId)
	if err != nil {
		return nil
	}
	// Call API
	resp, err := dms.CallAPI(
		"GET",
		"/workspace/"+workspaceId,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil
	}

	var workspace models.TwWorkspace
	err = json.Unmarshal(body, &workspace)
	if err != nil {
		return nil
	}

	return &workspace
}

func (s *WorkspaceService) DeleteWorkspace(id string) error {

	_, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("Invalid workspace id")
	}
	if id == "" {
		return errors.New("Invalid workspace id")
	}
	if a := s.GetWorkspaceById(id); a == nil {
		return errors.New("Workspace not found")
	}
	resp, err := dms.CallAPI(
		"DELETE",
		"/workspace/"+id,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Failed to delete workspace")
	}
	return nil

}
