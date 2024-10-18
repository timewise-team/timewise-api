package user_email

import (
	"api/dms"
	"encoding/json"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/user_email_dtos"
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
