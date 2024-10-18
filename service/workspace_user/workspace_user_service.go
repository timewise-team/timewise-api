package workspace_user

import (
	"api/dms"
	"encoding/json"
	"errors"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/workspace_user_dtos"
	"github.com/timewise-team/timewise-models/models"
	"io/ioutil"
	"net/http"
	"strconv"
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

func (s *WorkspaceUserService) GetWorkspaceUserList(workspaceID string) ([]workspace_user_dtos.GetWorkspaceUserListResponse, error) {
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

	var workspaceUserList []workspace_user_dtos.GetWorkspaceUserListResponse
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

func (s *WorkspaceUserService) DeleteWorkspaceUser(workspaceUser models.TwWorkspaceUser, workspaceUserMemberId string) error {
	workspaceID := workspaceUser.WorkspaceId
	var workspaceIDStr = strconv.Itoa(workspaceID)
	if workspaceUserMemberId == "" {
		return nil
	}

	// Call API
	resp, err := dms.CallAPI(
		"PUT",
		"/workspace_user/"+workspaceUserMemberId+"/workspace/"+workspaceIDStr,
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
		return err
	}

	WorkspaceLog := models.TwWorkspaceLog{
		WorkspaceId:     workspaceID,
		WorkspaceUserId: workspaceUser.ID,
		Action:          "delete",
		FieldChanged:    "workspace's user",
		OldValue:        workspaceUserMemberId,
		NewValue:        "",
		Description:     "Delete workspace's member",
	}

	err = s.AddWorkspaceLog(WorkspaceLog)

	return nil
}

func (s *WorkspaceUserService) AddWorkspaceLog(workspaceLog models.TwWorkspaceLog) error {

	// Call API
	resp, err := dms.CallAPI(
		"POST",
		"/workspace_log",
		workspaceLog,
		nil,
		nil,
		120,
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}

func (s *WorkspaceUserService) GetWorkspaceUserInvitationNotVerifiedList(workspaceID string) ([]workspace_user_dtos.GetWorkspaceUserListResponse, error) {
	// Call API
	resp, err := dms.CallAPI(
		"GET",
		"/workspace_user/invitation_not_verified/workspace/"+workspaceID,
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

func (s *WorkspaceUserService) UpdateWorkspaceUserRole(workspaceUser *models.TwWorkspaceUser, request workspace_user_dtos.UpdateWorkspaceUserRoleRequest) error {
	workspaceId := workspaceUser.WorkspaceId
	workspaceIdStr := strconv.Itoa(workspaceId)
	if workspaceIdStr == "" {
		return errors.New("workspace id not found")
	}
	// Call API
	resp, err := dms.CallAPI(
		"PUT",
		"/workspace_user/role/workspace/"+workspaceIdStr,
		request,
		nil,
		nil,
		120,
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err != nil || resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}
