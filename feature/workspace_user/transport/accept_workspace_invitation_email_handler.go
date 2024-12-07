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
		c.Set("Content-Type", "text/html")
		return c.Status(500).SendString(errorHtml("Failed to load config"))
	}
	token := c.Params("token")
	claims, err2 := auth_utils.ParseInvitationToken(token, cfg.JWT_SECRET)
	if err2 != nil {
		if errors.Is(err2, jwt.ErrTokenExpired) {
			workspaceId := claims["workspace_id"].(float64)
			workspaceIdStr := fmt.Sprintf("%.0f", workspaceId)
			workspaceUser, err3 := workspace_user.NewWorkspaceUserService().GetWorkspaceUserByEmailAndWorkspaceID(claims["email"].(string), workspaceIdStr)
			if err3 != nil || workspaceUser == nil {
				c.Set("Content-Type", "text/html")
				return c.Status(500).SendString(errorHtml("This request is invalid."))
			}
			if workspaceUser.Status == "pending" {
				err := workspace_user.NewWorkspaceUserService().UpdateStatusByEmailAndWorkspace(claims["email"].(string), workspaceId, "removed", false, true)
				if err != nil {
					c.Set("Content-Type", "text/html")
					return c.Status(fiber.StatusInternalServerError).SendString(errorHtml("Failed to update user status to 'removed': " + err.Error()))
				}
				c.Set("Content-Type", "text/html")
				return c.Status(fiber.StatusUnauthorized).SendString(errorHtml("Token expired. User status set to 'removed'."))
			}
		}
		c.Set("Content-Type", "text/html")
		return c.Status(fiber.StatusUnauthorized).SendString(errorHtml("This invitation link has been broken. Please request a new invitation."))
	}

	workspaceId := claims["workspace_id"].(float64)
	workspaceIdStr := fmt.Sprintf("%.0f", workspaceId)
	workspaceUser, err3 := workspace_user.NewWorkspaceUserService().GetWorkspaceUserByEmailAndWorkspaceID(claims["email"].(string), workspaceIdStr)
	if err3 != nil || workspaceUser == nil {
		c.Set("Content-Type", "text/html")
		return c.Status(500).SendString(errorHtml("This request is invalid."))
	}
	if workspaceUser.IsVerified && workspaceUser.Status == "joined" && workspaceUser.IsActive {
		c.Set("Content-Type", "text/html")
		return c.Status(400).SendString(errorHtml("User has already joined the workspace"))
	}

	if workspaceUser.Status != "joined" {
		isMember := claims["is_member"].(bool)
		var err error
		if isMember {
			err = workspace_user.NewWorkspaceUserService().UpdateStatusByEmailAndWorkspace(claims["email"].(string), workspaceId, "joined", true, false)
		} else {
			err = workspace_user.NewWorkspaceUserService().UpdateStatusByEmailAndWorkspace(claims["email"].(string), workspaceId, "joined", true, true)
		}
		if err != nil {
			c.Set("Content-Type", "text/html")
			return c.Status(500).SendString(errorHtml(err.Error()))
		}
		c.Set("Content-Type", "text/html")
		return c.SendString(successHtml("accept"))
	}
	c.Set("Content-Type", "text/html")
	return c.SendString(successHtml("already a member"))
}

func errorHtml(message string) string {
	return `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Invalid Email Request</title>
        <style>
            body {
                font-family: Arial, sans-serif;
                background-color: #f8f9fa;
                color: #343a40;
                display: flex;
                justify-content: center;
                align-items: center;
                height: 100vh;
                margin: 0;
            }
            .container {
                text-align: center;
                max-width: 500px;
                background: #fff;
                padding: 20px;
                border-radius: 10px;
                box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
            }
            h1 {
                color: #e74c3c;
                font-size: 24px;
            }
            p {
                font-size: 16px;
                line-height: 1.5;
                margin: 10px 0;
            }
            a {
                display: inline-block;
                margin-top: 20px;
                padding: 10px 20px;
                background-color: #007bff;
                color: white;
                text-decoration: none;
                border-radius: 5px;
            }
            a:hover {
                background-color: #0056b3;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <h1>Oops! Invalid Request</h1>
            <p>` + message + `</p>
            <p>If you think this is a mistake, please contact support for further assistance.</p>
        </div>
    </body>
    </html>
    `
}

func successHtml(action string) string {
	htmlContent := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Email ` + action + ` Success</title>
			<style>
				body { font-family: Arial, sans-serif; }
				.container { text-align: center; margin-top: 50px; }
				.success { color: green; font-size: 20px; }
				.error { color: red; font-size: 20px; }
				.button { padding: 10px 20px; font-size: 16px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 5px; }
			</style>
		</head>
		<body>
			<div class="container">
				<h1 class="success">Congratulations! Your email has been successfully ` + action + ` the workspace invitation.</h1>
				<p>If you ` + action + ` the workspace invitation, you has been added to this workspace accordingly.</p>`

	if action == "accept" {
		htmlContent += `
				<p>Your acceptance has been confirmed. You can now join the workspace.</p>`
	} else if action == "reject" {
		htmlContent += `
				<p>Your invitation has been rejected. If this was a mistake, please contact support.</p>`
	}
	htmlContent += `
				<a href="https://timewise.space/" class="button">You can close this page now.</a>
			</div>
		</body>
		</html>
	`
	return htmlContent
}
