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

// acceptInvitationViaEmail godoc
// @Summary Accept invitation via email
// @Description Accept invitation via email
// @Tags WorkspaceUser
// @Accept json
// @Produce json
// @Param token path string true "Token"
// @Success 200 {object} map[string]interface{} "Workspace invitation accepted successfully"
// @Failure 404 {object} map[string]interface{} "Workspace user not found"
// @Failure 401 {object} map[string]interface{} "Token expired or invalid"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/v1/workspace_user/accept-invitation-via-email/token/{token} [get]
func (h *WorkspaceUserHandler) acceptInvitationViaEmail(c *fiber.Ctx) error {
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
	if workspaceUser == nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "workspace user not found",
		})
	}
	if workspaceUser.IsVerified == true && workspaceUser.Status == "joined" && workspaceUser.IsActive == true {
		return c.Status(400).JSON(fiber.Map{
			"message": "User has already joined the workspace",
		})
	}
	if err2 != nil {
		if errors.Is(err2, jwt.ErrTokenExpired) {
			if workspaceUser.Status == "pending" {

				// Nếu token hết hạn, cập nhật trạng thái workspaceUser thành "removed".
				err := workspace_user.NewWorkspaceUserService().UpdateStatusByEmailAndWorkspace(claims["email"].(string), claims["workspace_id"].(float64), "removed", false, true)
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
	if workspaceUser.Status != "joined" {
		isMember := claims["is_member"].(bool)
		var err error
		if isMember {
			err = workspace_user.NewWorkspaceUserService().UpdateStatusByEmailAndWorkspace(claims["email"].(string), claims["workspace_id"].(float64), "joined", true, false)
		} else {
			err = workspace_user.NewWorkspaceUserService().UpdateStatusByEmailAndWorkspace(claims["email"].(string), claims["workspace_id"].(float64), "joined", true, true)
		}

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
			"message": "Workspace invitation accepted successfully",
		})
	}
	return c.JSON(fiber.Map{
		"message": "This user is already a member of this workspace",
	})
}
