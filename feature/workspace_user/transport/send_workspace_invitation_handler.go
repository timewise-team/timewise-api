package transport

import (
	"api/config"
	"api/notification"
	"api/service/auth"
	"api/service/user_email"
	"api/service/workspace"
	"api/service/workspace_user"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/workspace_user_dtos"
	"github.com/timewise-team/timewise-models/models"
	"strconv"
)

// sendInvitation godoc
// @Summary Send invitation to user (X-User-Email required, X-Workspace-Id required)
// @Description Send invitation to user (X-User-Email required, X-Workspace-Id required)
// @Tags WorkspaceUser
// @Accept json
// @Produce json
// @Security bearerToken
// @Param X-User-Email header string true "User Email"
// @Param X-Workspace-Id header string true "Workspace ID"
// @Param workspace_user body workspace_user_dtos.UpdateWorkspaceUserRoleRequest true "Workspace user object"
// @Success 200 {object} models.TwWorkspaceUser
// @Router /api/v1/workspace_user/send-invitation [post]
func (s *WorkspaceUserHandler) sendInvitation(c *fiber.Ctx) error {
	cfg, err1 := config.LoadConfig()
	if err1 != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to load config",
		})
	}
	workspaceUserLocal := c.Locals("workspace_user")
	if workspaceUserLocal == nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	workspaceUser, ok := workspaceUserLocal.(*models.TwWorkspaceUser)
	if !ok {
		return c.Status(400).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	var workspaceUserInvitationRequest workspace_user_dtos.UpdateWorkspaceUserRoleRequest
	if err := c.BodyParser(&workspaceUserInvitationRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}
	userEmail, errs := user_email.NewUserEmailService().GetUserEmail(workspaceUserInvitationRequest.Email)
	if userEmail == nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "This email is not registered",
		})
	}
	if errs != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// workspaceId, email
	// Check if any linked email is already in this workspace
	// If yes, then return error
	existingLinkedWorkspaceUser, checkErr := workspace_user.NewWorkspaceUserService().GetExistingLinkedWorkspaceUser(userEmail.Email, strconv.Itoa(workspaceUser.WorkspaceId))
	if checkErr != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	if len(existingLinkedWorkspaceUser) > 0 {
		return c.Status(500).JSON(fiber.Map{
			"message": "This email is already linked to another workspace",
			"email":   existingLinkedWorkspaceUser[0],
		})
	}

	workspaceUserCheck, err1 := workspace_user.NewWorkspaceUserService().GetWorkspaceUserByEmailAndWorkspaceID(userEmail.Email, strconv.Itoa(workspaceUser.WorkspaceId))
	if err1 != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	if workspaceUserCheck != nil {
		if workspaceUserCheck.Status == "joined" {
			return c.Status(200).JSON(fiber.Map{
				"message": "This user is already a member of this workspace",
			})
		}
		if workspaceUserCheck.Status == "pending" {
			worspaceIdStr := strconv.Itoa(workspaceUser.WorkspaceId)
			workspaceInfo := workspace.NewWorkspaceService().GetWorkspaceById(worspaceIdStr)
			acceptLink, declineLink, _ := auth.GenerateInviteLinks(cfg, userEmail.Email, workspaceInfo.ID, workspaceUserInvitationRequest.Role)
			content := auth.BuildInvitationContent(workspaceInfo, workspaceUserCheck.Role, acceptLink, declineLink)
			subject := fmt.Sprintf("Invitation to join workspace: %s", workspaceInfo.Title)
			if err := auth.SendInvitationEmail(cfg, userEmail.Email, content, subject); err != nil {
				return c.Status(500).JSON(fiber.Map{"message": "Failed to send invitation email"})
			}
			if err := PushInvitationNotification(workspaceInfo.Title, acceptLink, declineLink, userEmail.ID); err != nil {
				return c.Status(500).JSON(fiber.Map{"message": "Failed to send notification"})
			}
			return c.Status(200).JSON(fiber.Map{
				"message":      "Invitation sent successfully",
				"accept_link":  acceptLink,
				"decline_link": declineLink,
			})
		}
		if workspaceUserCheck.Status == "declined" || workspaceUserCheck.Status == "removed" {
			workspaceUserCheck.Status = "pending"
			workspaceUserCheck.Role = workspaceUserInvitationRequest.Role
			workspaceUserUpdate, err := workspace_user.NewWorkspaceUserService().UpdateWorkspaceUserStatus(workspaceUserCheck)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"message": "Internal server error",
				})
			}
			worspaceIdStr := strconv.Itoa(workspaceUser.WorkspaceId)
			workspaceInfo := workspace.NewWorkspaceService().GetWorkspaceById(worspaceIdStr)
			acceptLink, declineLink, _ := auth.GenerateInviteLinks(cfg, userEmail.Email, workspaceInfo.ID, workspaceUserInvitationRequest.Role)
			content := auth.BuildInvitationContent(workspaceInfo, workspaceUserUpdate.Role, acceptLink, declineLink)
			subject := fmt.Sprintf("Invitation to join workspace: %s", workspaceInfo.Title)
			if err := auth.SendInvitationEmail(cfg, userEmail.Email, content, subject); err != nil {
				return c.Status(500).JSON(fiber.Map{"message": "Failed to send invitation email"})
			}
			if err := PushInvitationNotification(workspaceInfo.Title, acceptLink, declineLink, userEmail.ID); err != nil {
				return c.Status(500).JSON(fiber.Map{"message": "Failed to send notification"})
			}
			return c.Status(200).JSON(fiber.Map{
				"message":      "Invitation sent successfully",
				"accept_link":  acceptLink,
				"decline_link": declineLink,
			})

		}

	}
	currentEmail := c.Locals("email")
	if currentEmail == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid session email"})
	}
	var workspaceUserResponse, err = workspace_user.NewWorkspaceUserService().AddWorkspaceUserInvitation(userEmail, workspaceUser.WorkspaceId, workspaceUserInvitationRequest, currentEmail.(string))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	worspaceIdStr := strconv.Itoa(workspaceUser.WorkspaceId)
	workspaceInfo := workspace.NewWorkspaceService().GetWorkspaceById(worspaceIdStr)
	acceptLink, declineLink, _ := auth.GenerateInviteLinks(cfg, userEmail.Email, workspaceInfo.ID, workspaceUserInvitationRequest.Role)
	content := auth.BuildInvitationContent(workspaceInfo, workspaceUserResponse.Role, acceptLink, declineLink)
	subject := fmt.Sprintf("Invitation to join workspace: %s", workspaceInfo.Title)
	if err := auth.SendInvitationEmail(cfg, userEmail.Email, content, subject); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to send invitation email"})
	}
	if err := PushInvitationNotification(workspaceInfo.Title, acceptLink, declineLink, userEmail.ID); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to send notification"})
	}
	return c.Status(200).JSON(fiber.Map{
		"message":      "Invitation sent successfully",
		"accept_link":  acceptLink,
		"decline_link": declineLink,
	})

}

// sendInvitationByMember godoc
// @Summary Send invitation to user by member (X-User-Email required, X-Workspace-Id required)
// @Description Send invitation to user by member (X-User-Email required, X-Workspace-Id required)
// @Tags WorkspaceUser
// @Accept json
// @Produce json
// @Security bearerToken
// @Param X-User-Email header string true "User Email"
// @Param X-Workspace-Id header string true "Workspace ID"
// @Param workspace_user body workspace_user_dtos.UpdateWorkspaceUserRoleRequest true "Workspace user object"
// @Success 200 {object} models.TwWorkspaceUser
// @Router /api/v1/workspace_user/member/send-invitation [post]
func (s *WorkspaceUserHandler) sendInvitationByMember(c *fiber.Ctx) error {
	cfg, err1 := config.LoadConfig()
	if err1 != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to load config",
		})
	}
	workspaceUserLocal := c.Locals("workspace_user")
	if workspaceUserLocal == nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	workspaceUser, ok := workspaceUserLocal.(*models.TwWorkspaceUser)
	if !ok {
		return c.Status(400).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	var workspaceUserInvitationRequest workspace_user_dtos.UpdateWorkspaceUserRoleRequest
	if err := c.BodyParser(&workspaceUserInvitationRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}
	userEmail, errs := user_email.NewUserEmailService().GetUserEmail(workspaceUserInvitationRequest.Email)
	if userEmail == nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "This email is not registered",
		})
	}
	if errs != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// workspaceId, email
	// Check if any linked email is already in this workspace
	// If yes, then return error
	existingLinkedWorkspaceUser, checkErr := workspace_user.NewWorkspaceUserService().GetExistingLinkedWorkspaceUser(userEmail.Email, strconv.Itoa(workspaceUser.WorkspaceId))
	if checkErr != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	if len(existingLinkedWorkspaceUser) > 0 {
		return c.Status(500).JSON(fiber.Map{
			"message": "This email is already linked to another workspace",
			"email":   existingLinkedWorkspaceUser[0],
		})
	}

	workspaceUserCheck, err1 := workspace_user.NewWorkspaceUserService().GetWorkspaceUserByEmailAndWorkspaceID(userEmail.Email, strconv.Itoa(workspaceUser.WorkspaceId))
	if err1 != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	if workspaceUserCheck != nil {
		if workspaceUserCheck.Status == "joined" {
			return c.Status(200).JSON(fiber.Map{
				"message": "This user is already a member of this workspace",
			})
		}
		if workspaceUserCheck.Status == "pending" {
			worspaceIdStr := strconv.Itoa(workspaceUser.WorkspaceId)
			workspaceInfo := workspace.NewWorkspaceService().GetWorkspaceById(worspaceIdStr)
			acceptLink, declineLink, _ := auth.GenerateInviteByMemberLinks(cfg, userEmail.Email, workspaceInfo.ID, workspaceUserInvitationRequest.Role)
			content := auth.BuildInvitationContent(workspaceInfo, workspaceUserCheck.Role, acceptLink, declineLink)
			subject := fmt.Sprintf("Invitation to join workspace: %s", workspaceInfo.Title)
			if err := auth.SendInvitationEmail(cfg, userEmail.Email, content, subject); err != nil {
				return c.Status(500).JSON(fiber.Map{"message": "Failed to send invitation email"})
			}
			if err := PushInvitationNotification(workspaceInfo.Title, acceptLink, declineLink, userEmail.ID); err != nil {
				return c.Status(500).JSON(fiber.Map{"message": "Failed to send notification"})
			}
			return c.Status(200).JSON(fiber.Map{
				"message":      "Invitation sent successfully",
				"accept_link":  acceptLink,
				"decline_link": declineLink,
			})
		}
		if workspaceUserCheck.Status == "declined" || workspaceUserCheck.Status == "removed" {
			workspaceUserCheck.Status = "pending"
			workspaceUserCheck.Role = workspaceUserInvitationRequest.Role
			workspaceUserUpdate, err := workspace_user.NewWorkspaceUserService().UpdateWorkspaceUserStatus(workspaceUserCheck)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"message": "Internal server error",
				})
			}
			worspaceIdStr := strconv.Itoa(workspaceUser.WorkspaceId)
			workspaceInfo := workspace.NewWorkspaceService().GetWorkspaceById(worspaceIdStr)
			acceptLink, declineLink, _ := auth.GenerateInviteByMemberLinks(cfg, userEmail.Email, workspaceInfo.ID, workspaceUserInvitationRequest.Role)
			content := auth.BuildInvitationContent(workspaceInfo, workspaceUserUpdate.Role, acceptLink, declineLink)
			subject := fmt.Sprintf("Invitation to join workspace: %s", workspaceInfo.Title)
			if err := auth.SendInvitationEmail(cfg, userEmail.Email, content, subject); err != nil {
				return c.Status(500).JSON(fiber.Map{"message": "Failed to send invitation email"})
			}
			if err := PushInvitationNotification(workspaceInfo.Title, acceptLink, declineLink, userEmail.ID); err != nil {
				return c.Status(500).JSON(fiber.Map{"message": "Failed to send notification"})
			}
			return c.Status(200).JSON(fiber.Map{
				"message":      "Invitation sent successfully",
				"accept_link":  acceptLink,
				"decline_link": declineLink,
			})

		}

	}
	currentEmail := c.Locals("email")
	if currentEmail == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid session email"})
	}
	var workspaceUserResponse, err = workspace_user.NewWorkspaceUserService().AddWorkspaceUserInvitation(userEmail, workspaceUser.WorkspaceId, workspaceUserInvitationRequest, currentEmail.(string))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	worspaceIdStr := strconv.Itoa(workspaceUser.WorkspaceId)
	workspaceInfo := workspace.NewWorkspaceService().GetWorkspaceById(worspaceIdStr)
	acceptLink, declineLink, _ := auth.GenerateInviteByMemberLinks(cfg, userEmail.Email, workspaceInfo.ID, workspaceUserInvitationRequest.Role)
	content := auth.BuildInvitationContent(workspaceInfo, workspaceUserResponse.Role, acceptLink, declineLink)
	subject := fmt.Sprintf("Invitation to join workspace: %s", workspaceInfo.Title)
	if err := auth.SendInvitationEmail(cfg, userEmail.Email, content, subject); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to send invitation email"})
	}
	if err := PushInvitationNotification(workspaceInfo.Title, acceptLink, declineLink, userEmail.ID); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to send notification"})
	}
	return c.Status(200).JSON(fiber.Map{
		"message":      "Invitation sent successfully",
		"accept_link":  acceptLink,
		"decline_link": declineLink,
	})

}

func PushInvitationNotification(workspaceTitle, acceptLink, declineLink string, userEmailId int) error {
	// create json of link
	link := map[string]string{
		"accept":  acceptLink,
		"decline": declineLink,
	}
	linkJson, _ := json.Marshal(link)

	// send notification
	notificationDto := models.TwNotifications{
		Title:       fmt.Sprintf("Invitation to join workspace %s", workspaceTitle),
		Description: fmt.Sprintf("You have been invited to join workspace %s", workspaceTitle),
		Link:        string(linkJson),
		UserEmailId: userEmailId,
		Type:        "workspace_invitation",
		Message:     "",
		IsSent:      true,
	}
	err := notification.PushNotifications(notificationDto)
	if err != nil {
		return err
	}
	return nil
}
