package transport

import (
	"api/config"
	"api/notification"
	"api/service/account"
	"api/service/auth"
	auth_utils "api/utils/auth"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
	"github.com/timewise-team/timewise-models/models"
	"strconv"
	"time"
)

type AccountHandler struct {
	service account.AccountService
}

func NewAccountHandler() *AccountHandler {
	service := account.NewAccountService()
	return &AccountHandler{
		service: *service,
	}
}

// getUserInfo godoc
// @Summary Get user info
// @Description Get user info
// @Tags account
// @Security bearerToken
// @Accept json
// @Produce json
// @Param status query string false "Status"
// @Success 200 {object} core_dtos.GetUserResponseDto
// @Router /api/v1/account/user [get]
func (h *AccountHandler) getUserInfo(c *fiber.Ctx) error {
	// get userId from context
	userId := c.Locals("userid")
	if userId == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid userId"})
	}
	status := c.Query("status")
	// call service to query database
	userInfo, err := h.service.GetUserInfoByUserId(userId.(string), status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// return user info
	return c.Status(fiber.StatusOK).JSON(userInfo)
}

// updateUserInfo godoc
// @Summary Update user info
// @Description Update user info
// @Tags account
// @Security bearerToken
// @Accept json
// @Produce json
// @Param updateUserInfoRequest body core_dtos.UpdateProfileRequestDto true "Update user info request"
// @Success 200 {object} core_dtos.GetUserResponseDto
// @Router /api/v1/account/user [patch]
func (h *AccountHandler) updateUserInfo(c *fiber.Ctx) error {
	// get userId from context
	userId := c.Locals("userid")
	if userId == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid userId"})
	}
	// get request body
	var updateUserInfoRequest core_dtos.UpdateProfileRequestDto
	if err := c.BodyParser(&updateUserInfoRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// call utils to validate request
	//if !account_utils.IsValidInputUpdateProfileRequest(updateUserInfoRequest) {
	//	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	//}
	// call service to update user info
	userResp, err := h.service.UpdateUserInfo(userId.(string), updateUserInfoRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(userResp)
}

// getLinkedUserEmails godoc
// @Summary Get linked user emails
// @Description Get linked user emails
// @Tags account
// @Security bearerToken
// @Param status query string false "Status"
// @Accept json
// @Produce json
// @Success 200 {array} string
// @Router /api/v1/account/user/emails [get]
func (h *AccountHandler) getLinkedUserEmails(c *fiber.Ctx) error {
	// get userId from context
	userId := c.Locals("userid")
	if userId == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid userId"})
	}
	status := c.Query("status")
	// call service to query database
	userEmails, err := h.service.GetLinkedUserEmails(userId.(string), status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// return user info
	return c.Status(fiber.StatusOK).JSON(userEmails)
}

// sendLinkEmailRequest godoc
// @Summary Send link an email request
// @Description Send link an email request
// @Tags account
// @Security bearerToken
// @Accept json
// @Produce json
// @Param email query string true "Target email"
// @Success 200 {object} string "Link email request sent"
// @Router /api/v1/account/user/emails/send [post]
func (h *AccountHandler) sendLinkEmailRequest(c *fiber.Ctx) error {
	// get userId from context
	userId := c.Locals("userid")
	if userId == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid userId"})
	}
	userIdStr, ok := userId.(string) // Kiểm tra xem userId có phải là kiểu string không
	if !ok {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID type")
	}

	userIdInt, err := strconv.Atoi(userIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID format")
	}
	// get email from request
	targetEmail := c.Query("email")
	if targetEmail == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email is required"})
	}
	// get current email from context
	currentEmail := c.Locals("email")
	if currentEmail == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid session email"})
	}
	// call service
	userEmailResp, err2 := h.service.SendLinkAnEmailRequest(userIdStr, targetEmail, currentEmail.(string))
	if err2 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err2.Error()})
	}
	// send notification
	notificationDto := models.TwNotifications{
		UserEmailId: userIdInt,
		Type:        "link email",
		Message: "A confirmation link has been successfully sent to " + targetEmail +
			". Please check your inbox and click the link to confirm your request. " +
			"Your request will be expired after 10 minutes." +
			"Please confirm of reject it before that time." +
			"If you don’t see the email, check your Spam or Promotions folder.",
		RelatedItemId:   userEmailResp.ID,
		RelatedItemType: "user_email",
		Title:           "Link Email Request",
		Description:     "A confirmation link has been successfully sent to " + targetEmail,
		IsSent:          true,
	}
	requestEmail, acceptLink, rejectLink, err := generateMessageEmail(userIdStr, targetEmail, userEmailResp.ExpiresAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	err = notification.PushNotifications(notificationDto)
	notificationDto = models.TwNotifications{
		UserEmailId:     userEmailResp.ID,
		Type:            "link email",
		Message:         requestEmail,
		RelatedItemId:   userIdInt,
		RelatedItemType: "user_email",
		Title:           "Link Email Request",
		Description:     "You have received a request to link your email address to account: " + currentEmail.(string),
		Link:            "Click here to approve: " + acceptLink + "<br>Click here to reject: " + rejectLink,
	}
	err = notification.PushNotifications(notificationDto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Link email request sent", "acceptLink": acceptLink})
}

func generateMessageEmail(userId string, email string, exprireAt *time.Time) (string, string, string, error) {
	exprireTime := "NULL"
	if exprireAt != nil {
		loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
		if err != nil {
			return "", "", "", fmt.Errorf("Error loading location: %v", err)
		}
		exprireAtTemp := exprireAt.In(loc)

		exprireTime = exprireAtTemp.Format("2006-01-02 15:04")
	}
	cfg, err1 := config.LoadConfig()
	if err1 != nil {
		return "", "", "", err1
	}

	accptLink, err := auth.GenerateLinkEmailLinks(cfg, userId, email, "linked")
	if err != nil {
		return "", "", "", err
	}

	rejectLink, err := auth.GenerateLinkEmailLinks(cfg, userId, email, "rejected")
	if err != nil {
		return "", "", "", err
	}

	// Nội dung email HTML
	emailContent := `
	<!DOCTYPE html>
	<html>
		<body>
			<p>Hello,</p>
			<p>You have requested to register the email address ` + email + `.</p>
			<p><strong>This request will be expired after 10 minutes. Please decide before: ` + exprireTime + `</strong></p>
			<p><strong>Please consider carefully, because these links cannot be use after once you clicked!!!</strong></p>
			<p>Please confirm or decline this request by clicking on one of the links below:</p>
			<p>
				<a href="` + accptLink + `">Confirm Registration</a><br>
				<a href="` + rejectLink + `">Decline Registration</a>
			</p>
			<p>If you did not request this registration, please ignore this message.</p>
			<p>Best regards,<br>Timewise Team</p>
		</body>
	</html>
`
	return emailContent, accptLink, rejectLink, nil
}

// actionEmailLinkRequest godoc
// @Summary Action email link request
// @Description Action email link request
// @Tags account
// @Accept json
// @Produce json
// @Param token path string true "Token"
// @Success 200 {object} core_dtos.GetUserResponseDto
// @Router /api/v1/account/user/emails/link/{token} [get]
func (h *AccountHandler) actionEmailLinkRequest(c *fiber.Ctx) error {
	cfg, err1 := config.LoadConfig()
	if err1 != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to load config",
		})
	}
	token := c.Params("token")
	claims, err2 := auth_utils.ParseInvitationToken(token, cfg.JWT_SECRET)
	if err2 != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token: " + err2.Error(),
		})
	}
	userId := claims["user_id"].(string)
	email := claims["email"].(string)
	action := claims["action"].(string)
	// call service to send mail
	_, err := h.service.UpdateStatusLinkEmailRequest(userId, email, action)
	if err != nil {
		if err.Error() == "email is not pending" {
			c.Set("Content-Type", "text/html")
			return c.SendString(errorHtml())
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	htmlContent := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Email ` + action + ` Success</title>
			<style>
				body { font-family: Arial, sans-serif; }
				.container { text-align: center; margin-top: 50px; }
				.success { color: green; font-size: 20px; }
				.error { color: red; font-size: 20px; }
				.button { padding: 10px 20px; font-size: 16px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 5px; }
			</style>
		</head>
		<body>
			<div class="container">
				<h1 class="success">Congratulations! Your email has been successfully ` + action + `.</h1>
				<p>If you ` + action + ` the email request, your account information has been updated accordingly.</p>`

	// If the action is "accept" or "reject", you can add specific messages or buttons.
	if action == "accept" {
		htmlContent += `
				<p>Your email registration has been confirmed. You can now access your account.</p>`
	} else if action == "reject" {
		htmlContent += `
				<p>Your email registration has been rejected. If this was a mistake, please contact support.</p>`
	}
	htmlContent += `
				<a href="/" class="button">You can close this page now.</a>
			</div>
		</body>
		</html>
	`
	// Send HTML content as response
	c.Set("Content-Type", "text/html")
	return c.SendString(htmlContent)
}

func errorHtml() string {
	return `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Invalid Email Request</title>
        <style>
            body {
                font-family: Arial, sans-serif;
                background-color: #f8f9fa;
                color: #343a40;
                display: flex;
                justify-content: center;
                align-items: center;
                height: 100vh;
                margin: 0;
            }
            .container {
                text-align: center;
                max-width: 500px;
                background: #fff;
                padding: 20px;
                border-radius: 10px;
                box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
            }
            h1 {
                color: #e74c3c;
                font-size: 24px;
            }
            p {
                font-size: 16px;
                line-height: 1.5;
                margin: 10px 0;
            }
            a {
                display: inline-block;
                margin-top: 20px;
                padding: 10px 20px;
                background-color: #007bff;
                color: white;
                text-decoration: none;
                border-radius: 5px;
            }
            a:hover {
                background-color: #0056b3;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <h1>Oops! Invalid Request</h1>
            <p>The email request you are trying to access is no longer valid or has already been processed.</p>
            <p>If you think this is a mistake, please contact support for further assistance.</p>
        </div>
    </body>
    </html>
    `
}

// unlinkAnEmail godoc
// @Summary Unlink an email
// @Description Unlink an email
// @Tags account
// @Security bearerToken
// @Accept json
// @Produce json
// @Param email query string true "Email"
// @Success 200 {object} core_dtos.GetUserResponseDto
// @Router /api/v1/account/user/emails/unlink [post]
func (h *AccountHandler) unlinkAnEmail(c *fiber.Ctx) error {
	targetEmail := c.Query("email")
	if targetEmail == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email is required"})
	}
	// call service
	userResp, err := h.service.UnlinkAnEmail(targetEmail)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// get userId from context
	userId := c.Locals("userid")
	if userId == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid userId"})
	}
	userIdStr, ok := userId.(string) // Kiểm tra xem userId có phải là kiểu string không
	if !ok {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID type")
	}

	userIdInt, err := strconv.Atoi(userIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID format")
	}
	// send notification
	notificationDto := models.TwNotifications{
		Title:       "Email Removal",
		Description: "Email " + targetEmail + " has been removed successfully",
		UserEmailId: userIdInt,
		Type:        "unlink email",
		Message:     "Unlink to email: " + targetEmail + " successfully",
	}
	err = notification.PushNotifications(notificationDto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	notificationDto = models.TwNotifications{
		UserEmailId: userResp.ID,
		Type:        "unlink email",
		Message:     "You have been unlinked from: " + targetEmail,
	}
	err = notification.PushNotifications(notificationDto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(userResp)
}

// deactivateAccount godoc
// @Summary Deactivate account
// @Description Deactivate account
// @Tags account
// @Security bearerToken
// @Accept json
// @Produce json
// @Success 200 {object} string "Account deactivated"
// @Failure 400 {object} fiber.Map "Invalid userId"
// @Failure 500 {object} fiber.Map "Internal server error"
// @Router /api/v1/account/user/deactivate [post]
func (h *AccountHandler) deactivateAccount(c *fiber.Ctx) error {
	// get userId from context
	userId := c.Locals("userid")
	if userId == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid userId"})
	}
	// call service to deactivate account
	err := h.service.DeactivateAccount(userId.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Account deactivated"})
}

// getParentLinkedEmails godoc
// @Summary Get parent linked emails
// @Description Get parent linked emails
// @Tags account
// @Security bearerToken
// @Accept json
// @Produce json
// @Success 200 {array} string
// @Router /api/v1/account/user/emails/parent [get]
func (h *AccountHandler) getParentLinkedEmails(c *fiber.Ctx) error {
	// get userId from context
	currentEmail := c.Locals("email")
	if currentEmail == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid userId"})
	}
	// call service to query database
	userEmails, err := h.service.GetParentLinkedEmails(currentEmail.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	// return parent email
	return c.Status(fiber.StatusOK).JSON(userEmails)
}

// clearStatusRejectedEmail godoc
// @Summary Clear status rejected email
// @Description Clear status rejected email
// @Tags account
// @Security bearerToken
// @Accept json
// @Produce json
// @Param email query string true "Email"
// @Success 200 {object} string "Clear status rejected email successfully"
// @Router /api/v1/account/user/emails/clear-rejected [get]
func (h *AccountHandler) clearStatusRejectedEmail(c *fiber.Ctx) error {
	// get target email from request
	targetEmail := c.Query("email")
	if targetEmail == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email is required"})
	}
	// call service
	err := h.service.ClearStatusRejectedEmail(targetEmail)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Clear status rejected email successfully"})
}
