package workspace_user

import (
	"api/dms"
	"api/service/workspace"
	auth_utils "api/utils/auth"
	"encoding/json"
	"errors"
	"fmt"
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

	//Validate workspace request
	if email == "" {
		return nil, errors.New("Email is required")
	}
	if workspaceID == "" {
		return nil, errors.New("Workspace ID is required")
	}
	_, err := strconv.Atoi(workspaceID)
	if err != nil {
		return nil, errors.New("Invalid workspace ID")
	}
	if auth_utils.IsValidEmail(email) == false {
		return nil, errors.New("Invalid email")
	}
	if _, err := strconv.Atoi(workspaceID); err != nil {
		return nil, errors.New("Invalid workspace ID")
	}

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
	if workspaceID == "" {
		return nil, errors.New("workspace id not found")
	}
	if _, err := strconv.Atoi(workspaceID); err != nil {
		return nil, errors.New("workspace id is invalid")
	}
	if workspace.NewWorkspaceService().GetWorkspaceById(workspaceID).ID == 0 {
		return nil, errors.New("workspace not found")
	}
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

func (s *WorkspaceUserService) GetWorkspaceUserListForManage(workspaceID string) ([]workspace_user_dtos.GetWorkspaceUserListResponse, error) {
	// Call API
	resp, err := dms.CallAPI(
		"GET",
		"/workspace_user/manage/workspace/"+workspaceID,
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
	if workspaceID == "" {
		return nil, errors.New("workspace id not found")
	}
	if _, err := strconv.Atoi(workspaceID); err != nil {
		return nil, errors.New("workspace id is invalid")
	}
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

func (s *WorkspaceUserService) DeleteWorkspaceUser(workspaceUser *models.TwWorkspaceUser, workspaceUserMemberId string) error {
	workspaceID := workspaceUser.WorkspaceId
	var workspaceIDStr = strconv.Itoa(workspaceID)
	if workspaceUserMemberId == "" {
		return nil
	}
	_, err := strconv.Atoi(workspaceUserMemberId)
	if err != nil {
		return errors.New("workspace user is unvalid")
	}

	// Call API
	resp, err := dms.CallAPI(
		"DELETE",
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

	workspaceUserInfo, err := s.GetWorkspaceUserInformation(workspaceUserMemberId)
	if err != nil {
		return err
	}
	if workspaceUserInfo.ID == 0 {
		return errors.New("workspace user not found")
	}

	WorkspaceLog := models.TwWorkspaceLog{
		WorkspaceId:     workspaceID,
		WorkspaceUserId: workspaceUser.ID,
		Action:          "delete",
		FieldChanged:    "workspace's user",
		OldValue:        workspaceUserInfo.LastName + " " + workspaceUserInfo.FirstName + " - " + workspaceUserInfo.Email,
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
	if workspaceID == "" {
		return nil, errors.New("workspace id not found")
	}
	if _, err := strconv.Atoi(workspaceID); err != nil {
		return nil, errors.New("workspace id is invalid")
	}
	if workspace.NewWorkspaceService().GetWorkspaceById(workspaceID).ID == 0 {
		return nil, errors.New("workspace not found")
	}
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

func (s *WorkspaceUserService) VerifyWorkspaceUserInvitation(workspaceUser *models.TwWorkspaceUser, email string) error {
	if email == "" {
		return errors.New("email not found")
	}
	if !auth_utils.IsValidEmail(email) {
		return errors.New("email is invalid")
	}
	workspaceId := workspaceUser.WorkspaceId
	workspaceIdStr := strconv.Itoa(workspaceId)
	if workspaceIdStr == "" {
		return errors.New("workspace id not found")
	}
	if workspace.NewWorkspaceService().GetWorkspaceById(workspaceIdStr).ID == 0 {
		return errors.New("workspace not found")
	}
	// Call API
	resp, err := dms.CallAPI(
		"PUT",
		"/workspace_user/verify-invitation/workspace/"+workspaceIdStr+"/email/"+email,
		nil,
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
func (s *WorkspaceUserService) DisproveWorkspaceUserInvitation(workspaceUser *models.TwWorkspaceUser, email string) error {
	workspaceId := workspaceUser.WorkspaceId
	workspaceIdStr := strconv.Itoa(workspaceId)
	if workspaceIdStr == "" {
		return errors.New("workspace id not found")
	}
	if workspace.NewWorkspaceService().GetWorkspaceById(workspaceIdStr).ID == 0 {
		return errors.New("workspace not found")
	}
	// Call API
	resp, err := dms.CallAPI(
		"PUT",
		"/workspace_user/disprove-invitation/workspace/"+workspaceIdStr+"/email/"+email,
		nil,
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

func (s *WorkspaceUserService) AddWorkspaceUserInvitation(userEmail *models.TwUserEmail, workspaceID int, request workspace_user_dtos.UpdateWorkspaceUserRoleRequest) (models.TwWorkspaceUser, error) {
	if userEmail.ID == 0 {
		return models.TwWorkspaceUser{}, errors.New("user email id not found")
	}

	var workspaceUser = models.TwWorkspaceUser{
		UserEmailId: userEmail.ID,
		WorkspaceId: workspaceID,
		Role:        request.Role,
		Status:      "pending",
		IsActive:    false,
		IsVerified:  true,
	}
	workspaceIDStr := strconv.Itoa(workspaceID)
	if workspaceIDStr == "" {
		return models.TwWorkspaceUser{}, errors.New("workspace id not found")
	}
	// Call API
	resp, err := dms.CallAPI(
		"POST",
		"/workspace_user",
		workspaceUser,
		nil,
		nil,
		120,
	)
	if err != nil {
		return models.TwWorkspaceUser{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return models.TwWorkspaceUser{}, err
	}
	if len(body) == 0 || string(body) == "null" {
		return models.TwWorkspaceUser{}, errors.New("response is null or empty")
	}

	var workspaceUserResponse models.TwWorkspaceUser
	err = json.Unmarshal(body, &workspaceUserResponse)
	if err != nil {
		return models.TwWorkspaceUser{}, err
	}

	return workspaceUserResponse, nil
}

func (s *WorkspaceUserService) AddWorkspaceUserViaScheduleInvitation(userEmail *models.TwUserEmail, workspaceID int, isVerified bool) (models.TwWorkspaceUser, error) {
	var workspaceUser = models.TwWorkspaceUser{
		UserEmailId: userEmail.ID,
		WorkspaceId: workspaceID,
		Role:        "guest",
		Status:      "pending",
		IsActive:    false,
		IsVerified:  isVerified,
	}
	workspaceIDStr := strconv.Itoa(workspaceID)
	if workspaceIDStr == "" {
		return models.TwWorkspaceUser{}, errors.New("workspace id not found")
	}
	// Call API
	resp, err := dms.CallAPI(
		"POST",
		"/workspace_user",
		workspaceUser,
		nil,
		nil,
		120,
	)
	if err != nil {
		return models.TwWorkspaceUser{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return models.TwWorkspaceUser{}, err
	}
	if len(body) == 0 || string(body) == "null" {
		return models.TwWorkspaceUser{}, errors.New("response is null or empty")
	}

	var workspaceUserResponse models.TwWorkspaceUser
	err = json.Unmarshal(body, &workspaceUserResponse)
	if err != nil {
		return models.TwWorkspaceUser{}, err
	}

	return workspaceUserResponse, nil
}

func (s *WorkspaceUserService) UpdateWorkspaceUserStatus(check *models.TwWorkspaceUser) (*models.TwWorkspaceUser, error) {
	// Call API
	resp, err := dms.CallAPI(
		"PUT",
		"/workspace_user/update-status/"+strconv.Itoa(check.ID),
		check,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error")
	}

	return check, nil

}

func (s *WorkspaceUserService) GetWorkspaceUserInformation(workspaceUserId string) (workspace_user_dtos.GetWorkspaceUserListResponse, error) {
	if workspaceUserId == "" {
		return workspace_user_dtos.GetWorkspaceUserListResponse{}, errors.New("workspace user id not found")
	}
	if _, err := strconv.Atoi(workspaceUserId); err != nil {
		return workspace_user_dtos.GetWorkspaceUserListResponse{}, errors.New("workspace user id is invalid")
	}
	// Call API
	resp, err := dms.CallAPI(
		"GET",
		"/workspace_user/"+workspaceUserId+"/info",
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return workspace_user_dtos.GetWorkspaceUserListResponse{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return workspace_user_dtos.GetWorkspaceUserListResponse{}, err
	}

	var workspaceUser workspace_user_dtos.GetWorkspaceUserListResponse
	err = json.Unmarshal(body, &workspaceUser)
	if err != nil {
		return workspace_user_dtos.GetWorkspaceUserListResponse{}, err
	}

	return workspaceUser, nil
}

func (s *WorkspaceUserService) AcceptInvitation(workspaceUser *models.TwWorkspaceUser, workspaceID int) error {
	workspaceIDStr := strconv.Itoa(workspaceID)
	if workspaceIDStr == "" {
		return errors.New("workspace id not found")
	}
	// Call API
	resp, err := dms.CallAPI(
		"PUT",
		"/workspace_user/accept-invitation/"+strconv.Itoa(workspaceUser.ID)+"/workspace/"+workspaceIDStr,
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
		return errors.New("error")
	}

	return nil
}

func (s *WorkspaceUserService) UpdateStatusByEmailAndWorkspace(email string, workspaceID float64, status string, isActive bool) error {
	workspaceIDStr := fmt.Sprintf("%.0f", workspaceID)
	if workspaceIDStr == "" {
		return errors.New("workspace id not found")
	}
	// Call API
	resp, err := dms.CallAPI(
		"PUT",
		"/workspace_user/update-status/email/"+email+"/workspace/"+workspaceIDStr+"/status/"+status+"/is_active/"+strconv.FormatBool(isActive),
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
		return errors.New("error")
	}

	return nil
}
