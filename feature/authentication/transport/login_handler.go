package transport

import (
	"api/feature/authentication/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/user_login_dtos"
)

// Login @Summary Login
// @Description Login
// @Tags auth
// @Accept json
// @Produce json
// @Param body body user_login_dtos.UserLoginRequest true "User login request"
// @Success 200 {object} user_login_dtos.UserLoginResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req user_login_dtos.UserLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}
	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username and password are required",
		})
	}
	userResponse, err := usecase.Login(req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})

	}
	return c.JSON(userResponse)
}
