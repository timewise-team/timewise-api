package workspace_log

import (
	"api/dms"
	"api/service/workspace"
	"encoding/json"
	"errors"
	"github.com/timewise-team/timewise-models/models"
	"io/ioutil"
	"net/http"
	"strconv"
)

type WorkspaceLogService struct {
	// Thêm các dependencies cần thiết nếu có (ví dụ: database, API client, v.v.)
}

func NewWorkspaceLogService() *WorkspaceLogService {
	return &WorkspaceLogService{}
}

func (s *WorkspaceLogService) GetWorkspaceLogs(workspaceId string) ([]models.TwWorkspaceLog, error) {
	// Validate workspaceId
	if workspaceId == "" {
		return nil, errors.New("Invalid workspaceId")
	}
	if _, err := strconv.Atoi(workspaceId); err != nil {
		return nil, errors.New("Invalid workspaceId")
	}
	if workspace.NewWorkspaceService().GetWorkspaceById(workspaceId) == nil {
		return nil, errors.New("Workspace not found")
	}

	// Call API
	resp, err := dms.CallAPI(
		"GET",
		"/workspace_log/workspace/"+workspaceId,
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

	var workspaceLogs []models.TwWorkspaceLog
	err = json.Unmarshal(body, &workspaceLogs)
	if err != nil {
		return nil, err
	}

	return workspaceLogs, nil
}
