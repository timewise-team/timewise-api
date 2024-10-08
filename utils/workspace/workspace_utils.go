package workspace

import (
	auth_utils "api/utils/auth"
	"errors"
	dtos "github.com/timewise-team/timewise-models/dtos/core_dtos/create_workspace_dtos"
	"net/url"
)

func ValidateWorkspace(workspace dtos.CreateWorkspaceRequest) error {
	if workspace == (dtos.CreateWorkspaceRequest{}) {
		return errors.New("workspace is required")
	}
	if workspace.Title == "" {
		return errors.New("workspace title is required")
	}
	if len(workspace.Title) > 50 {
		return errors.New("workspace title must not exceed 100 characters")
	}
	if workspace.Description == "" {
		return errors.New("workspace description is required")
	}
	if len(workspace.Description) > 500 {
		return errors.New("workspace description must not exceed 500 characters")
	}
	unescapedEmail, err := url.QueryUnescape(workspace.Email)
	if err != nil {
		return errors.New("invalid email format")
	}
	if !auth_utils.IsValidEmail(unescapedEmail) {
		return errors.New("invalid email format")
	}
	return nil
}
