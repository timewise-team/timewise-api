package middleware

import (
	auth_service "api/service/auth"
	"api/service/workspace_user"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func CheckWorkspaceRole(requiredRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		email := c.Get("X-User-Email")
		if email == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "X-User-Email header is required",
			})
		}

		// Lấy workspace_id từ body hoặc header
		workspaceID := c.Get("X-Workspace-ID")
		if workspaceID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "X-Workspace-ID header is required",
			})
		}
		//Check email in email list
		var userId = c.Locals("userid")
		if userId == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Access denied",
			})
		}
		if userId == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Access denied",
			})
		}
		checkEmail, err := auth_service.NewAuthService().CheckEmailInList(userId.(string), email)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}
		if !checkEmail {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Access denied",
			})
		}
		// Lấy vai trò của người dùng trong workspace
		workspaceUser, err := workspace_user.NewWorkspaceUserService().GetWorkspaceUserByEmailAndWorkspaceID(email, workspaceID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}

		if workspaceUser == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Access denied",
			})
		}

		// Kiểm tra xem vai trò có phù hợp không
		hasRole := false
		for _, role := range requiredRoles {
			if strings.ToLower(workspaceUser.Role) == strings.ToLower(role) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Forbidden",
			})
		}

		// Lưu thông tin workspaceUser vào context để sử dụng sau này nếu cần
		c.Locals("workspace_user", workspaceUser)

		return c.Next()
	}
}
