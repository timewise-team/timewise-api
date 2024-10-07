package transport

import (
	"api/dms"
	auth_service "api/service/auth"
	auth_utils "api/utils/auth"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	dtos "github.com/timewise-team/timewise-models/dtos/core_dtos/user_register_dtos"
	"github.com/timewise-team/timewise-models/models"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type GoogleAuthRequest struct {
	Credentials string `json:"credentials"`
}

type GoogleAuthResponse struct {
	AccessToken      string               `json:"access_token"`
	ExpiresIn        int                  `json:"expires_in"`
	TokenType        string               `json:"token_type"`
	IsNewUser        bool                 `json:"is_new_user"`
	IdToken          string               `json:"id_token"`
	LinkedUserEmails []models.TwUserEmail `json:"linked_user_emails"`
}
type GetUserEmailSyncResponse []models.TwUserEmail

// @Summary Google callback
// @Description Google callback
// @Tags auth
// @Accept json
// @Produce json
// @Param body body GoogleAuthRequest true "Google auth request"
// @Success 200 {object} GoogleAuthResponse
// @Router /api/v1/auth/callback [post]
func (h *AuthHandler) googleCallback(c *fiber.Ctx) error {
	var req GoogleAuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request",
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

	getOrCreateUserReq := dtos.GetOrCreateUserRequestDto{
		Email:          oauthData.Email,
		FullName:       oauthData.Name,
		ProfilePicture: oauthData.Picture,
		VerifiedEmail:  oauthData.VerifiedEmail,
		GoogleId:       oauthData.Id,
		GivenName:      oauthData.GivenName,
		FamilyName:     oauthData.FamilyName,
		Locale:         oauthData.Locale,
	}

	// Get or Create user
	resp, err := dms.CallAPI(
		"POST",
		"/user/get-create",
		getOrCreateUserReq,
		nil,
		nil,
		120,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not get or create user"})
	}

	// marshal response body
	var userRespDto dtos.GetOrCreateUserResponseDto
	err = json.Unmarshal(body, &userRespDto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not marshal response body"})
	}

	if userRespDto.IsNewUser {
		_, err := auth_service.NewAuthService().InitNewUser(userRespDto.User)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not init new user"})
		}
	}
	idStr := strconv.Itoa(userRespDto.User.ID)
	resp2, err := dms.CallAPI(
		"GET",
		"/user_email/user/"+idStr,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body2, err := ioutil.ReadAll(resp2.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not get or create user"})
	}

	// marshal response body
	var userEmailSync GetUserEmailSyncResponse
	err = json.Unmarshal(body2, &userEmailSync)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not marshal response body"})
	}
	// Generate JWT token
	accessToken, expiresIn, err := auth_utils.GenerateJWTToken(userRespDto.User, viper.GetString("JWT_SECRET"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate JWT token"})
	}

	// Send the token back to the frontend
	return c.JSON(GoogleAuthResponse{
		AccessToken:      accessToken,
		ExpiresIn:        expiresIn,
		TokenType:        "Bearer",
		IsNewUser:        userRespDto.IsNewUser,
		IdToken:          req.Credentials,
		LinkedUserEmails: userEmailSync,
	})
}
