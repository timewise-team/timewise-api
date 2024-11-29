package auth

import (
	"api/config"
	"api/dms"
	"api/service/notification_setting"
	auth_utils "api/utils/auth"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/timewise-team/timewise-models/models"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"log"
	"net/http"
)

type GetUserEmailSyncResponse []models.TwUserEmail

type AuthService struct {
	// Thêm các dependencies cần thiết nếu có (ví dụ: database, API client, v.v.)
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) CreateUserEmail(user models.TwUser) (*models.TwUserEmail, error) {
	var UserEmail models.TwUserEmail
	UserEmail.Email = user.Email
	UserEmail.UserId = user.ID
	UserEmail.User = user

	resp, err := dms.CallAPI(
		"POST",
		"/user_email",
		UserEmail,
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

	var UserEmailResponse models.TwUserEmail
	err = json.Unmarshal(body, &UserEmailResponse)
	if err != nil {
		return nil, err
	}

	return &UserEmailResponse, nil
}

// CreateWorkspace handles creating the workspace via API
func (s *AuthService) CreateWorkspace() (*models.TwWorkspace, error) {
	var WorkspaceRequest models.TwWorkspace
	WorkspaceRequest.Title = "personal"
	WorkspaceRequest.Type = "personal"

	resp, err := dms.CallAPI(
		"POST",
		"/workspace",
		WorkspaceRequest,
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

	var WorkspaceResponse models.TwWorkspace
	err = json.Unmarshal(body, &WorkspaceResponse)
	if err != nil {
		return nil, err
	}

	return &WorkspaceResponse, nil
}
func (s *AuthService) CreateWorkspaceUser(userEmail *models.TwUserEmail, workspace *models.TwWorkspace) (*models.TwWorkspaceUser, error) {
	var WorkspaceUserRequest models.TwWorkspaceUser
	WorkspaceUserRequest.UserEmailId = userEmail.ID
	WorkspaceUserRequest.WorkspaceId = workspace.ID
	WorkspaceUserRequest.Workspace = *workspace
	WorkspaceUserRequest.UserEmail = *userEmail
	WorkspaceUserRequest.Role = "owner"
	WorkspaceUserRequest.Status = "joined"
	WorkspaceUserRequest.IsActive = true
	WorkspaceUserRequest.IsVerified = true
	WorkspaceUserRequest.ExtraData = ""
	WorkspaceUserRequest.WorkspaceKey = ""

	resp, err := dms.CallAPI(
		"POST",
		"/workspace_user",
		WorkspaceUserRequest,
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

	var WorkspaceUserResponse models.TwWorkspaceUser
	err = json.Unmarshal(body, &WorkspaceUserResponse)
	if err != nil {
		return nil, err
	}

	return &WorkspaceUserResponse, nil

}

func (s *AuthService) InitNewUser(user models.TwUser) (bool, error) {
	//_ = s.CreateNotificationSetting(user)
	// Create user email
	userEmailResponse, err := s.CreateUserEmail(user)
	if err != nil {
		return false, err // Return error if email creation fails
	}

	// Create workspace
	workspaceResponse, err := s.CreateWorkspace()
	if err != nil {
		return false, err // Return error if workspace creation fails
	}

	_, err = s.CreateWorkspaceUser(userEmailResponse, workspaceResponse)
	// Create workspace user
	if err != nil {
		return false, err // Return error if workspace user creation fails
	}
	err = notification_setting.NewNotificationSettingService().CreateNotificationSetting(user.ID)
	if err != nil {
		return false, err
	}
	return true, nil // Success
}

func (s *AuthService) CheckEmailInList(userId string, email string) (bool, error) {
	resp, err := dms.CallAPI(
		"GET",
		"/user_email/user/"+userId,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Println(resp.StatusCode)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false, err
	}
	var userEmailSync GetUserEmailSyncResponse
	err = json.Unmarshal(body, &userEmailSync)
	if err != nil {
		return false, err
	}
	for _, userEmail := range userEmailSync {
		if userEmail.Email == email {
			return true, nil
		}
	}
	return false, nil
}

//func (s *AuthService) CreateNotificationSetting(user models.TwUser) *models.TwNotificationSettings {
//	var NotificationSetting models.TwNotificationSettings
//	NotificationSetting.UserId = user.ID
//	NotificationSetting.NotificationOnComment = true
//	NotificationSetting.NotificationOnDueDate = true
//	NotificationSetting.NotificationOnScheduleChange = true
//	NotificationSetting.NotificationOnDueDate = true
//
//	resp, err := dms.CallAPI(
//		"POST",
//		"/notification_setting",
//		NotificationSetting,
//		nil,
//		nil,
//		120,
//	)
//	if err != nil {
//		return nil
//	}
//
//	defer resp.Body.Close()
//
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil || resp.StatusCode != http.StatusOK {
//		return nil
//	}
//
//	var NotificationSettingResponse models.TwNotificationSettings
//	err = json.Unmarshal(body, &NotificationSettingResponse)
//	if err != nil {
//		return nil
//	}
//
//	return &NotificationSettingResponse
//
//}

// SendEmail thực hiện gửi email với cấu hình và nội dung đã tạo.
func SendEmail(dialer *gomail.Dialer, message *gomail.Message) error {
	if err := dialer.DialAndSend(message); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}
	log.Println("Email sent successfully")
	return nil
}

func ConfigSMTP(cfg *config.Config) *gomail.Dialer {
	return gomail.NewDialer(cfg.SMPHost, cfg.SMTPPort, cfg.SMTPEmail, cfg.SMTPPassword)
}

func GenerateInviteLinks(cfg *config.Config, email string, workspaceId int, role string) (string, string, error) {
	acceptToken, err := auth_utils.GenerateInvitationToken(workspaceId, "accept", cfg.JWT_SECRET, email, role)
	if err != nil {
		return "", "", err
	}

	declineToken, err := auth_utils.GenerateInvitationToken(workspaceId, "decline", cfg.JWT_SECRET, email, role)
	if err != nil {
		return "", "", err
	}

	acceptLink := fmt.Sprintf("%s/workspace_user/accept-invitation-via-email/token/%s", cfg.BaseURL, acceptToken)
	declineLink := fmt.Sprintf("%s/workspace_user/decline-invitation-via-email/token/%s", cfg.BaseURL, declineToken)

	return acceptLink, declineLink, nil
}

func GenerateInviteScheduleLinks(cfg *config.Config, scheduleId int, workspaceUserId int) (string, string, error) {
	acceptToken, err := auth_utils.GenerateScheduleInvitationToken(workspaceUserId, "accept", cfg.JWT_SECRET, scheduleId)
	if err != nil {
		return "", "", err
	}

	declineToken, err := auth_utils.GenerateScheduleInvitationToken(workspaceUserId, "decline", cfg.JWT_SECRET, scheduleId)
	if err != nil {
		return "", "", err
	}

	acceptLink := fmt.Sprintf("%s/schedule_participant/accept-invitation-via-email/token/%s", cfg.BaseURL, acceptToken)
	declineLink := fmt.Sprintf("%s/schedule_participant/decline-invitation-via-email/token/%s", cfg.BaseURL, declineToken)

	return acceptLink, declineLink, nil
}

func GenerateLinkEmailLinks(cfg *config.Config, currentUid string, email string, action string) (string, error) {
	token, err := auth_utils.GenerateLinkEmailToken(currentUid, email, action, cfg.JWT_SECRET)
	if err != nil {
		return "", err
	}
	link := fmt.Sprintf("%s/account/user/emails/link/%s", cfg.BaseURL, token)

	return link, nil
}

func BuildInvitationContent(info *models.TwWorkspace, role, acceptLink, declineLink string) string {
	return fmt.Sprintf(`
	<html>
		<head>
			<style>
				body {
					font-family: 'Arial', sans-serif;
					background-color: #f5f6fa;
					color: #333;
					line-height: 1.6;
					margin: 0;
					padding: 20px;
				}
				.container {
					max-width: 600px;
					margin: 0 auto;
					background-color: white;
					box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
					border-radius: 10px;
					overflow: hidden;
				}
				.header {
					background-color: #4a90e2;
					color: white;
					padding: 20px;
					text-align: center;
					font-size: 24px;
					font-weight: bold;
				}
				.content {
					padding: 20px;
				}
				.btn {
					display: inline-block;
					margin-top: 10px;
					padding: 12px 30px;
					border-radius: 5px;
					text-decoration: none;
					font-weight: bold;
					color: white;
					transition: background-color 0.3s ease;
				}
				.btn-accept {
					background-color: #28a745;
				}
				.btn-accept:hover {
					background-color: #218838;
				}
				.btn-decline {
					background-color: #dc3545;
					margin-left: 10px;
				}
				.btn-decline:hover {
					background-color: #c82333;
				}
				.footer {
					margin-top: 20px;
					font-size: 14px;
					color: #999;
					text-align: center;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">Workspace Invitation</div>
				<div class="content">
					<p>Hello,</p>
					<p>You have been invited to join the workspace: <b>%s</b>.</p>
					<p>Your role: <b>%s</b></p>
					<a href="%s" class="btn btn-accept">Accept Invitation</a>
					<a href="%s" class="btn btn-decline">Decline Invitation</a>
				</div>
				<div class="footer">
					<p>If you have any questions, feel free to contact our support team.</p>
				</div>
			</div>
		</body>
	</html>
	`, info.Title, role, acceptLink, declineLink)
}

func BuildScheduleInvitationContent(info *models.TwSchedule, acceptLink, declineLink string) string {
	return fmt.Sprintf(`
	<html>
		<head>
			<style>
				body {
					font-family: 'Arial', sans-serif;
					background-color: #f5f6fa;
					color: #333;
					line-height: 1.6;
					margin: 0;
					padding: 20px;
				}
				.container {
					max-width: 600px;
					margin: 0 auto;
					background-color: white;
					box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
					border-radius: 10px;
					overflow: hidden;
				}
				.header {
					background-color: #4a90e2;
					color: white;
					padding: 20px;
					text-align: center;
					font-size: 24px;
					font-weight: bold;
				}
				.content {
					padding: 20px;
				}
				.btn {
					display: inline-block;
					margin-top: 10px;
					padding: 12px 30px;
					border-radius: 5px;
					text-decoration: none;
					font-weight: bold;
					color: white;
					transition: background-color 0.3s ease;
				}
				.btn-accept {
					background-color: #28a745;
				}
				.btn-accept:hover {
					background-color: #218838;
				}
				.btn-decline {
					background-color: #dc3545;
					margin-left: 10px;
				}
				.btn-decline:hover {
					background-color: #c82333;
				}
				.footer {
					margin-top: 20px;
					font-size: 14px;
					color: #999;
					text-align: center;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">Schedule Invitation</div>
				<div class="content">
					<p>Hello,</p>
					<p>You have been invited to join the schedule: <b>%s</b>.</p>
					<a href="%s" class="btn btn-accept">Accept Invitation</a>
					<a href="%s" class="btn btn-decline">Decline Invitation</a>
				</div>
				<div class="footer">
					<p>If you have any questions, feel free to contact our support team.</p>
				</div>
			</div>
		</body>
	</html>
	`, info.Title, acceptLink, declineLink)
}

func SendInvitationEmail(cfg *config.Config, email string, content string, subject string) error {
	for _, smtpConfig := range smtpConfigs {
		// Cấu hình SMTP
		dialer := gomail.NewDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.Email, smtpConfig.Password)
		if dialer == nil {
			log.Printf("Failed to configure SMTP dialer for %s", smtpConfig.Email)
			continue
		}

		// Tạo message mới
		message := gomail.NewMessage()
		message.SetHeader("From", smtpConfig.Email)
		message.SetHeader("To", email)
		message.SetHeader("Subject", subject)
		message.SetBody("text/html", content)

		// Gửi email
		if err := dialer.DialAndSend(message); err != nil {
			log.Printf("Failed to send email using %s: %v", smtpConfig.Email, err)
			continue
		}

		log.Println("Invitation email sent successfully")
		return nil
	}

	return errors.New("failed to send invitation email with all SMTP configurations")
}

var smtpConfigs = []struct {
	Host     string
	Port     int
	Email    string
	Password string
}{
	{"smtp.gmail.com", 587, "timewise.space@gmail.com", "dczt wlvd eisn cixf"},
	{"smtp.gmail.com", 587, "khanhhnhe170088@fpt.edu.vn", "cddn ujge aqlm xmjb"},
	{"smtp.gmail.com", 587, "khanhhn.hoang@gmail.com", "dgbx xyvw ciqg txbl"},
	{"smtp.gmail.com", 587, "khanhhnhe170088@fpt.edu.vn", "iaqw vmoj fxgb zzne"},
	{"smtp.gmail.com", 587, "thuandqhe170881@fpt.edu.vn", "whzq ivlb hevo jhdi"},
	{"smtp.gmail.com", 587, "builanviet@gmail.com", "lowo laid zgda chnc"},
	{"smtp.gmail.com", 587, "ngkkhanh006@gmail.com", "soet mdxg doio fmrt"},
}
