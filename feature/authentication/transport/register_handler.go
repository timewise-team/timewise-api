package transport

import (
	"api/feature/authentication/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/user_register_dtos"
)

// @Summary Register
// @Description Register
// @Tags auth
// @Accept json
// @Produce json
// @Param body body user_register_dto.RegisterRequestDto true "User register request"
// @Success 200 {object} user_register_dto.RegisterResponseDto
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) register(c *fiber.Ctx) error {
	var RegisterRequestDto user_register_dto.RegisterRequestDto
	if err := c.BodyParser(&RegisterRequestDto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := usecase.RegisterUser(RegisterRequestDto)
	if err != nil {
		if err.Error() == "passwords do not match" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Register successfully",
	})
}
