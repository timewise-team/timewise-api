package transport

import (
	"api/service/account"
	auth_utils "api/utils/auth"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
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
// @Success 200 {object} core_dtos.GetUserResponseDto
// @Router /api/v1/account/user [get]
func (h *AccountHandler) getUserInfo(c *fiber.Ctx) error {
	// get userId from context
	userId := c.Locals("userid")
	if userId == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid userId"})
	}
	// call service to query database
	userInfo, err := h.service.GetUserInfoByUserId(userId.(string))
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
	// call service to query database
	userEmails, err := h.service.GetLinkedUserEmails(userId.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// return user info
	return c.Status(fiber.StatusOK).JSON(userEmails)
}

// linkAnEmail godoc
// @Summary Link an email
// @Description Link an email
// @Tags account
// @Security bearerToken
// @Accept json
// @Produce json
// @Param linkAnEmailRequest body core_dtos.GoogleAuthRequest true "Link an email request"
// @Success 200 {object} core_dtos.GetUserResponseDto
// @Router /api/v1/account/user/emails [post]
func (h *AccountHandler) linkAnEmail(c *fiber.Ctx) error {
	// get userId from context
	userId := c.Locals("userid")
	if userId == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid userId"})
	}
	var req core_dtos.GoogleAuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request",
		})
	}
	if req.Credentials == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Credentials is required",
		})
	}
	// decode credentials
	decodedCredentials, err := auth_utils.VerifyGoogleToken(req.Credentials)
	var oauthData auth_utils.GoogleOauthData
	err = json.Unmarshal(decodedCredentials, &oauthData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not decode credentials",
		})
	}
	// call service
	userResp, err := h.service.LinkAnEmail(userId.(string), oauthData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(userResp)
}

// unlinkAnEmail godoc
// @Summary Unlink an email
// @Description Unlink an email
// @Tags account
// @Security bearerToken
// @Accept json
// @Produce json
// @Param unlinkAnEmailRequest body core_dtos.GoogleAuthRequest true "Unlink an email request"
// @Success 200 {object} core_dtos.GetUserResponseDto
// @Router /api/v1/account/user/emails/unlink [post]
func (h *AccountHandler) unlinkAnEmail(c *fiber.Ctx) error {
	var req core_dtos.GoogleAuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request",
		})
	}
	if req.Credentials == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Credentials is required",
		})
	}
	// decode credentials
	decodedCredentials, err := auth_utils.VerifyGoogleToken(req.Credentials)
	var oauthData auth_utils.GoogleOauthData
	err = json.Unmarshal(decodedCredentials, &oauthData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not decode credentials",
		})
	}
	// call service
	userResp, err := h.service.UnlinkAnEmail(oauthData.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(userResp)
}
