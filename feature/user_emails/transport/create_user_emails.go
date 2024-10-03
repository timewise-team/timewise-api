package transport

import (
	"api/dms"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"io/ioutil"
	"net/http"
	"time"
)

type CreateUserEmailsRequest struct {
	UserId int    `json:"user_id"`
	Email  string `json:"email"`
}
type CreateUserEmailsResponse struct {
	ID        int           `gorm:"primary_key"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	DeletedAt time.Time     `json:"deleted_at" gorm:"default:null"`
	UserId    int           `json:"user_id" gorm:"index"`
	Email     string        `json:"email"`
	User      models.TwUser `gorm:"foreignkey:UserId;association_foreignkey:ID"`
}

// createNewUserEmail godoc
// @Summary Create new user email
// @Description Create new user email
// @Tags User Emails
// @Accept json
// @Produce json
// @Success 200 {object} CreateUserEmailsResponse
// @Router /api/v1/user_email [post]
func (h *UserEmailsHandler) createNewUserEmail(c *fiber.Ctx) error {
	var req CreateUserEmailsRequest
	userId, ok := c.Locals("userId").(int)
	if !ok {
		return c.Status(400).JSON(fiber.Map{
			"message": "userId not found",
		})
	}
	email, ok := c.Locals("email").(string)
	if !ok {
		return c.Status(400).JSON(fiber.Map{
			"message": "email not found",
		})
	}
	req.UserId = userId
	req.Email = email
	resp, err := dms.CallAPI(
		"POST",
		"/user_email",
		req,
		nil,
		nil,
		120,
	)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not get or create user"})
	}

	var data CreateUserEmailsResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not marshal response body"})
	}

	return c.Status(200).JSON(data)
}
