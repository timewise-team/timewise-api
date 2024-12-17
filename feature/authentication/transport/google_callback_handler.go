package transport

import (
	"api/dms"
	"api/service/account"
	auth_service "api/service/auth"
	auth_utils "api/utils/auth"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
	dtos "github.com/timewise-team/timewise-models/dtos/core_dtos/user_register_dtos"
	"github.com/timewise-team/timewise-models/models"
	"io/ioutil"
	"net/http"
)

type GetUserEmailSyncResponse []models.TwUserEmail

// @Summary Google callback
// @Description Google callback
// @Tags auth
// @Accept json
// @Produce json
// @Param body body core_dtos.GoogleAuthRequest true "Google auth request"
// @Success 200 {object} core_dtos.GoogleAuthResponse
// @Router /api/v1/auth/callback [post]
func (h *AuthHandler) googleCallback(c *fiber.Ctx) error {
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
	checkLinkedEmail, err := account.NewAccountService().GetParentLinkedEmails(oauthData.Email)
	if err != nil {
		print("No checkLinkedEmail", err)
		//return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not get linked email"})
	}
	if len(checkLinkedEmail) > 0 {
		User, err := account.NewAccountService().GetUserByEmail(checkLinkedEmail)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not get user by email"})
		}
		if User.ID != 0 {
			getOrCreateUserReq = dtos.GetOrCreateUserRequestDto{
				Email:          User.Email,
				FullName:       User.LastName + " " + User.FirstName,
				ProfilePicture: User.ProfilePicture,
				VerifiedEmail:  User.IsVerified,
				GoogleId:       User.GoogleId,
				GivenName:      User.FirstName,
				FamilyName:     User.LastName,
				Locale:         User.Locale,
			}
		}
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

	// Generate JWT token
	accessToken, expiresIn, err := auth_utils.GenerateJWTToken(userRespDto.User, viper.GetString("JWT_SECRET"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate JWT token"})
	}

	// Send the token back to the frontend
	return c.JSON(core_dtos.GoogleAuthResponse{
		AccessToken: accessToken,
		ExpiresIn:   expiresIn,
		TokenType:   "Bearer",
		IsNewUser:   userRespDto.IsNewUser,
		IdToken:     req.Credentials,
	})
}
