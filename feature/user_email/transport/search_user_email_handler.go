package transport

import (
	"api/service/user_email"
	"github.com/gofiber/fiber/v2"
)

// search user email
// @Summary search user email
// @Description search user email
// @Tags User Email
// @Accept json
// @Produce json
// @Param query path string true "query"
// @Success 200 {array} user_email_dtos.SearchUserEmailResponse
// @Router /api/v1/user_email/search-user_email/{query} [get]
func (h *UserEmailHandler) searchUserEmail(c *fiber.Ctx) error {
	query := c.Params("query")

	var userEmail, err = user_email.NewUserEmailService().SearchUserEmail(query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if userEmail == nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to search user email",
		})
	}
	return c.JSON(userEmail)
}
