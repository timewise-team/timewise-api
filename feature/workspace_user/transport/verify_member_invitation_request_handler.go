package transport

import "github.com/gofiber/fiber/v2"

func (h *WorkspaceUserHandler) verifyMemberInvitationRequest(c *fiber.Ctx) error {
	//var workspaceUserRequest workspace_user_dtos.TwWorkspaceUserRequest
	//if err := c.BodyParser(&workspaceUserRequest); err != nil {
	//	return c.Status(400).JSON(fiber.Map{
	//		"message": err.Error(),
	//	})
	//}
	//
	//var workspaceUser, err = workspace_user.NewWorkspaceUserService().VerifyMemberInvitationRequest(workspaceUserRequest)
	//if err != nil {
	//	return c.Status(500).JSON(fiber.Map{
	//		"message": err.Error(),
	//	})
	//}
	//if workspaceUser == nil {
	//	return c.Status(500).JSON(fiber.Map{
	//		"message": "Failed to verify member invitation request",
	//	})
	//}
	//return c.JSON(workspaceUser)
	return nil
}
