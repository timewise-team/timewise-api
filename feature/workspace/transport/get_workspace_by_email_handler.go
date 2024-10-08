package transport

import (
	"api/service/workspace"
	auth_utils "api/utils/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"net/url"
)

type GetWorkspaceByEmailResponse struct {
	Workspaces []models.TwWorkspace `json:"workspaces"`
}

// getWorkspacesByEmail godoc
// @Summary Get workspaces by email
// @Description Get workspaces by email
// @Tags workspace
// @Accept json
// @Produce json
// @Param email path string false "Email"
// @Success 200 {array} models.TwWorkspace
// @Router /api/v1/workspace/get-workspaces-by-email/{email} [get]
func (h *WorkspaceHandler) getWorkspacesByEmail(c *fiber.Ctx) error {

	email := c.Params("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}
	emailFix, err1 := url.QueryUnescape(email)
	if err1 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err1.Error(),
		})
	}
	if email != "all" && !auth_utils.IsValidEmail(emailFix) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email format",
		})
	}
	userId := c.Locals("userid")
	if userId == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user not found",
		})
	}
	if userId == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user not found",
		})
	}
	userIdStr, ok := userId.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error parsing user id",
		})
	}
	var (
		workspaces []models.TwWorkspace
		err        error
	)
	if email == "all" {
		workspaces, err = workspace.NewWorkspaceService().GetWorkspacesByUserId(userIdStr)
	} else {
		workspaces, err = workspace.NewWorkspaceService().GetWorkspacesByEmail(emailFix)
	}
	// Call service
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if len(workspaces) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No workspaces found",
		})
	}
	// Return the response

	return c.JSON(workspaces)
}
