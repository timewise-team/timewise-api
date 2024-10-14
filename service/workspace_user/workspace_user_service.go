package workspace_user

import (
	"api/dms"
	"encoding/json"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/workspace_user_dtos"
	"github.com/timewise-team/timewise-models/models"
	"io/ioutil"
	"net/http"
)

type WorkspaceUserService struct {
	// Thêm các dependencies cần thiết nếu có (ví dụ: database, API client, v.v.)
}

func NewWorkspaceUserService() *WorkspaceUserService {
	return &WorkspaceUserService{}
}

//func (s *WorkspaceUserService) CheckEmail(email string,userId string) (bool, error) {
//	resp, err := dms.CallAPI(
//		"GET",
//		"/workspace_user/email/"+email+"/user_id/"+userId,
//		nil,
//		nil,
//		nil,
//		120,
//	)
//	if err != nil {
//		return false, err
//	}
//	defer resp.Body.Close()
//	if resp.StatusCode != http.StatusOK {
//		return false, errors.New("error")
//	}
//	return true, nil
//
//}

func (s *WorkspaceUserService) GetWorkspaceUserByEmailAndWorkspaceID(email string, workspaceID string) (*models.TwWorkspaceUser, error) {
	//userId := c.Locals("userid")
	//if userId == nil {
	//	return nil, errors.New("user not found")
	//}
	//if userId == "" {
	//	return nil, errors.New("user not found")
	//}
	//userIdStr, ok := userId.(string)
	//if !ok {
	//	return nil, errors.New("error parsing user id")
	//}

	//if !ok {
	//	return nil, errors.New("error parsing user id")
	//}

	// Call API
	resp, err := dms.CallAPI(
		"GET",
		"/workspace_user/email/"+email+"/workspace/"+workspaceID,
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

	var workspaceUser models.TwWorkspaceUser
	err = json.Unmarshal(body, &workspaceUser)
	if err != nil {
		return nil, err
	}

	return &workspaceUser, nil
}

func (s *WorkspaceUserService) GetWorkspaceUserList(workspaceID string) ([]models.TwWorkspaceUser, error) {
	// Call API
	resp, err := dms.CallAPI(
		"GET",
		"/workspace_user/workspace/"+workspaceID,
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

	var workspaceUserList []models.TwWorkspaceUser
	err = json.Unmarshal(body, &workspaceUserList)
	if err != nil {
		return nil, err
	}

	return workspaceUserList, nil
}

func (s *WorkspaceUserService) GetWorkspaceUserInvitationList(workspaceID string) ([]workspace_user_dtos.GetWorkspaceUserListResponse, error) {
	// Call API
	resp, err := dms.CallAPI(
		"GET",
		"/workspace_user/invitation/workspace/"+workspaceID,
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

	var workspaceUserList []workspace_user_dtos.GetWorkspaceUserListResponse
	err = json.Unmarshal(body, &workspaceUserList)
	if err != nil {
		return nil, err
	}

	return workspaceUserList, nil
}
