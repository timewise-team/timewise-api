package workspace

import (
	"api/dms"
	"encoding/json"
	"errors"
	"github.com/timewise-team/timewise-models/models"
	"io/ioutil"
	"net/http"
)

type WorkspaceService struct {
	// Thêm các dependencies cần thiết nếu có (ví dụ: database, API client, v.v.)
}

func NewWorkspaceService() *WorkspaceService {
	return &WorkspaceService{}
}

func (s *WorkspaceService) GetWorkspacesByEmail(email string) ([]models.TwWorkspace, error) {
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
