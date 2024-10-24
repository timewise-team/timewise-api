package transport

import (
	"api/config"
	"api/service/workspace_user"
	auth_utils "api/utils/auth"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// declineInvitationViaEmail godoc
// @Summary Decline invitation via email
// @Description Decline invitation via email
// @Tags WorkspaceUser
// @Accept json
// @Produce json
// @Param token path string true "Token"
// @Success 200 {object} map[string]interface{} "Invitation declined successfully"
// @Failure 404 {object} map[string]interface{} "Workspace user not found"
// @Failure 401 {object} map[string]interface{} "Invalid or expired token"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/v1/workspace_user/decline-invitation-via-email/token/{token} [get]
func (h *WorkspaceUserHandler) declineInvitationViaEmail(c *fiber.Ctx) error {
	cfg, err1 := config.LoadConfig()
	if err1 != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to load config",
		})
	}
	token := c.Params("token")
	claims, err2 := auth_utils.ParseInvitationToken(token, cfg.JWT_SECRET)
	workspaceId := claims["workspace_id"].(float64)

	workspaceIdStr := fmt.Sprintf("%.0f", workspaceId)
	workspaceUser, err3 := workspace_user.NewWorkspaceUserService().GetWorkspaceUserByEmailAndWorkspaceID(claims["email"].(string), workspaceIdStr)
	if err3 != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err3.Error(),
		})
	}
	if workspaceUser == nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "workspace user not found",
		})
	}
	if err2 != nil {
		if errors.Is(err2, jwt.ErrTokenExpired) {
			if workspaceUser.Status == "pending" {

				// Nếu token hết hạn, cập nhật trạng thái workspaceUser thành "removed".
				err := workspace_user.NewWorkspaceUserService().UpdateStatusByEmailAndWorkspace(claims["email"].(string), claims["workspace_id"].(float64), "removed", false)
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"message": "Failed to update user status to 'removed': " + err.Error(),
					})
				}
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Token expired. User status set to 'removed'.",
				})
			}
		}
		// Xử lý các lỗi khác liên quan đến token.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token: " + err2.Error(),
		})
	}
	var err = workspace_user.NewWorkspaceUserService().UpdateStatusByEmailAndWorkspace(claims["email"].(string), claims["workspace_id"].(float64), "declined", false)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if workspaceUser == nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to accept workspace invitation",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Workspace invitation declined successfully",
	})
}
