package transport

import (
	"api/service/workspace"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"log"
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
	emailFix, err1 := url.QueryUnescape(email)
	if err1 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err1.Error(),
		})
	}
	userId := c.Locals("userid").(string)
	var (
		workspaces []models.TwWorkspace
		err        error
	)
	log.Println("email: ", email)
	if email == "all" {
		workspaces, err = workspace.NewWorkspaceService().GetWorkspacesByUserId(userId)
	} else {
		workspaces, err = workspace.NewWorkspaceService().GetWorkspacesByEmail(emailFix)
	}
	// Call service
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(workspaces)
}
