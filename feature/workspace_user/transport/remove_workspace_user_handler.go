package transport

import (
	"github.com/gofiber/fiber/v2"
)

func (h *WorkspaceUserHandler) deleteWorkspaceUser(c *fiber.Ctx) error {
	//workspaceUserLocal := c.Locals("workspace_user")
	//if workspaceUserLocal == nil {
	//	return c.Status(400).JSON(fiber.Map{
	//		"message": "Access denied",
	//	})
	//}
	//workspaceUser, ok := workspaceUserLocal.(*models.TwWorkspaceUser)
	//if !ok {
	//	return c.Status(400).JSON(fiber.Map{
	//		"message": "Access denied",
	//	})
	//}
	//workspaceUserMember := c.Params("workspace_user_id")
	//if workspaceUserMember == "" {
	//	return c.Status(400).JSON(fiber.Map{
	//		"message": "member is required",
	//	})
	//}
	//err := workspace_user.NewWorkspaceUserService().DeleteWorkspaceUser(workspaceUser, workspaceUserMember)
	//if err != nil {
	//	return c.Status(500).JSON(fiber.Map{
	//		"message": err.Error(),
	//	})
	//}
	//
	//return c.JSON(fiber.Map{
	//	"message": "Workspace member deleted successfully",
	//})
	return nil
}
