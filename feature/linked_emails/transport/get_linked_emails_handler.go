package transport

import (
	"api/dms"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
	"github.com/timewise-team/timewise-models/models"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
// @Router /api/v1/linked_emails/ [get]
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
// @Param email query string true "Email"
// @Success 200 {array} models.TwUserEmail
// @Router /api/v1/linked_emails/ [post]
func (h *LinkedEmailsHandler) linkAnEmail(c *fiber.Ctx) error {
	// get email from params
	encodedEmail := c.Query("email")
	if encodedEmail == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email is required"})
	}
	email, err := url.QueryUnescape(encodedEmail)
	queryParam := map[string]string{
		"email": email,
	}
	// check if email is already a user
	resp, err := dms.CallAPI(
		"GET",
		"/user/",
		nil,
		nil,
		queryParam,
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
			var resquestDto = core_dtos.LinkEmailRequestDto{
				UserId: userIdStr,
				Email:  email,
			}
			// if user exists, link email to user (by change user_id in user_email table)
			_, err = dms.CallAPI(
				"POST",
				"/link-email/",
				resquestDto,
				nil,
				nil,
				120,
			)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
		} else {
			// if user does not exist, throw error
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "This email is not a user"})
		}
	}
	allEmails := h.getLinkedUserEmail(c)
	// return all linked email
	return c.Status(fiber.StatusOK).JSON(allEmails)
}
