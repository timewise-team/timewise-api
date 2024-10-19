package transport

import (
	"github.com/gofiber/fiber/v2"
)

func (s *WorkspaceUserHandler) sendInvitation(c *fiber.Ctx) error {
	//workspaceUserLocal := c.Locals("workspace_user")
	//if workspaceUserLocal == nil {
	//	return c.Status(400).JSON(fiber.Map{
	//		"message": "Access denied",
	//	})
	//}
	//workspaceUser, ok := workspaceUserLocal.(models.TwWorkspaceUser)
	//if !ok {
	//	return c.Status(400).JSON(fiber.Map{
	//		"message": "Access denied",
	//	})
	//}
	//workspaceUserStr := strconv.Itoa(workspaceUser.WorkspaceId)
	//var workspaceUserList, err = workspace_user.NewWorkspaceUserService().GetWorkspaceUserInvitationNotVerifiedList(workspaceUserStr)
	//if err != nil {
	//	return c.Status(500).JSON(fiber.Map{
	//		"message": err.Error(),
	//	})
	//}
	//if workspaceUserList == nil {
	//	return c.Status(500).JSON(fiber.Map{
	//		"message": "Failed to get workspace user list",
	//	})
	//}
	//return c.JSON(workspaceUserList)
	return nil
}
