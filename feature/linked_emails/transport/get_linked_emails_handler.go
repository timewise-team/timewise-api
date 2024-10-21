package transport

import (
	"api/dms"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"io/ioutil"
	"log"
	"net/http"
)

type GetUserEmailSyncResponse []models.TwUserEmail

// getLinkedUserEmail godoc
// @Summary Get linked user email
// @Description Get linked user email
// @Tags linked_emails
// @Security bearerToken
// @Accept json
// @Produce json
// @Success 200 {array} models.TwUserEmail
// @Router /api/v1/user-emails/get-linked-email [get]
func (h *LinkedEmailsHandler) getLinkedUserEmail(c *fiber.Ctx) error {
	userId := c.Locals("userid")
	if userId == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user not found",
		})
	}
	if userId == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user not found",
		})
	}
	userIdStr, ok := userId.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error parsing user id",
		})
	}
	resp, err := dms.CallAPI(
		"GET",
		"/user_email/user/"+userIdStr,
		nil,
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
	var userEmailSync GetUserEmailSyncResponse
	err = json.Unmarshal(body, &userEmailSync)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not marshal response body"})
	}
	if len(userEmailSync) == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "There are no emails linked to this user"})
	}
	return c.JSON(userEmailSync)
}

func (h *LinkedEmailsHandler) linkAnEmail(c *fiber.Ctx) error {
	return nil
}
