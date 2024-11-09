package user_email

import (
	"api/dms"
	"encoding/json"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/user_email_dtos"
	"github.com/timewise-team/timewise-models/models"
	"log"
)

type UserEmailService struct {
}

func NewUserEmailService() *UserEmailService {
	return &UserEmailService{}
}

func (service *UserEmailService) SearchUserEmail(query string) ([]user_email_dtos.SearchUserEmailResponse, error) {
	resp, err := dms.CallAPI(
		"GET",
		"/user_email/search/"+query,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, nil
	}
	defer resp.Body.Close()
	var userEmails []user_email_dtos.SearchUserEmailResponse
	if err := json.NewDecoder(resp.Body).Decode(&userEmails); err != nil {
		return nil, err
	}
	return userEmails, nil
}

func (service *UserEmailService) GetUserEmail(email string) (*models.TwUserEmail, error) {
	log.Println(email)
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
	var userEmail models.TwUserEmail
	if err := json.NewDecoder(resp.Body).Decode(&userEmail); err != nil {
		return nil, err
	}
	log.Printf("User email: %v", userEmail)
	return &userEmail, nil
}

func (service *UserEmailService) GetUserEmailInProgress(scheduleId string) (*[]user_email_dtos.UserEmailStatusResponse, error) {

	resp, err := dms.CallAPI(
		"GET",
		"/user_email/listApprove/"+scheduleId,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var userEmail []user_email_dtos.UserEmailStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&userEmail); err != nil {
		return nil, err
	}
	return &userEmail, nil
}
