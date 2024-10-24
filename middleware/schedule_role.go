package middleware

import (
	"api/service/schedule"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/schedule_participant_dtos"
	"github.com/timewise-team/timewise-models/models"
	"strconv"
	"strings"
)

// Middleware kiểm tra trạng thái của lịch và quyền của người dùng workspace.
func CheckScheduleStatus(requiredRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scheduleID := c.Params("scheduleId")
		if scheduleID == "" {
			var dto schedule_participant_dtos.InviteToScheduleRequest

			if err := c.BodyParser(&dto); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "Invalid request body",
				})
			}
			scheduleID = strconv.Itoa(dto.ScheduleId)
		}

		workspaceUser, ok := c.Locals("workspace_user").(*models.TwWorkspaceUser)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}

		workspaceUserIDStr := strconv.Itoa(workspaceUser.ID)

		scheduleParticipant, err := schedule.NewScheduleService().FetchScheduleParticipant(workspaceUserIDStr, scheduleID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}

		if workspaceUser.Role == "admin" || workspaceUser.Role == "owner" {
			c.Locals("scheduleParticipant", scheduleParticipant)
			return c.Next()
		} else {
			hasRole := false
			for _, role := range requiredRoles {
				if strings.ToLower(scheduleParticipant.Status) == strings.ToLower(role) {
					hasRole = true
					break
				}
			}

			if !hasRole {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"message": "Forbidden",
				})
			}

			c.Locals("scheduleParticipant", scheduleParticipant)
			return c.Next()
		}
	}
}
