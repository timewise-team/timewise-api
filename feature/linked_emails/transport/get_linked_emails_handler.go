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

// linkAnEmail godoc
// @Summary Link an email to a user
// @Description Link an email to a user
// @Tags linked_emails
// @Security bearerToken
// @Accept json
// @Produce json
// @Param email path string true "Email"
// @Success 200 {array} models.TwUserEmail
// @Router /api/v1/user-emails/link-email/{email} [post]
func (h *LinkedEmailsHandler) linkAnEmail(c *fiber.Ctx) error {
	// get email from params
	email := c.Params("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email is required"})
	}
	// check if email is already a user
	resp, err := dms.CallAPI(
		"GET",
		"/users/"+email,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.Body == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not get user"})
	} else {
		if resp.StatusCode == http.StatusOK {
			var user models.TwUser
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not read response body"})
			}
			err = json.Unmarshal(body, &user)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not marshal response body"})
			}
			var userEmail models.TwUserEmail
			userEmail.Email = email
			userEmail.UserId = user.ID
			userEmail.User = user
			// if user exists, link email to user
			_, err = dms.CallAPI(
				"POST",
				"/user-email/",
				userEmail,
				nil,
				nil,
				120,
			)
		} else {
			// if user does not exist, throw error
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "This email is not a user"})
		}
	}
	allEmails := h.getLinkedUserEmail(c)
	// return all linked email
	return c.Status(fiber.StatusOK).JSON(allEmails)
}
