package transport

import (
	"api/dms"
	auth_utils "api/utils/auth"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	dtos "github.com/timewise-team/timewise-models/dtos/core_dtos/user_register_dtos"
	"io/ioutil"
	"net/http"
)

type GoogleAuthRequest struct {
	Credentials string `json:"credentials"`
}

type GoogleAuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	IsNewUser   bool   `json:"is_new_user"`
	IdToken     string `json:"id_token"`
}

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

	// Generate JWT token
	accessToken, expiresIn, err := auth_utils.GenerateJWTToken(userRespDto.User, viper.GetString("JWT_SECRET"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate JWT token"})
	}

	// Send the token back to the frontend
	return c.JSON(GoogleAuthResponse{
		AccessToken: accessToken,
		ExpiresIn:   expiresIn,
		TokenType:   "Bearer",
		IsNewUser:   userRespDto.IsNewUser,
		IdToken:     req.Credentials,
	})
}
