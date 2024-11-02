package transport

import (
	"api/dms"
	workspace_utils "api/utils/workspace"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"io/ioutil"
	"net/http"
	"strconv"
)

type UpdateWorkspace struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// updateWorkspace godoc
// @Summary Update a workspace (X-User-Email required, X-Workspace-Id required)
// @Description Update a workspace (X-User-Email required, X-Workspace-Id required)
// @Tags workspace
// @Accept json
// @Produce json
// @Security bearerToken
// @Param X-User-Email header string true "User Email"
// @Param X-Workspace-Id header string true "Workspace ID"
// @Param workspace body UpdateWorkspace true "Workspace"
// @Success 200 {object} models.TwWorkspace
// @Router /api/v1/workspace/update-workspace [put]
func (h *WorkspaceHandler) updateWorkspace(c *fiber.Ctx) error {
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
	workspaceIdStr := strconv.Itoa(workspaceUser.WorkspaceId)
	if workspaceIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid workspaceId"})
	}
	var updateWorkspace UpdateWorkspace
	err := c.BodyParser(&updateWorkspace)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	response, err := dms.CallAPI(
		"GET",
		"/workspace/"+workspaceIdStr,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer response.Body.Close()
	var Workspace models.TwWorkspace
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if response.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get workspace"})
	}
	err = json.Unmarshal(body, &Workspace)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	err1 := workspace_utils.ValidateWorkspaces(updateWorkspace.Title, updateWorkspace.Description)
	if err1 != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err1.Error()})
	}
	Workspace.Title = updateWorkspace.Title
	Workspace.Description = updateWorkspace.Description

	resp, err := dms.CallAPI(
		"PUT",
		"/workspace/"+workspaceIdStr,
		Workspace,
		nil,
		nil,
		120,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	defer resp.Body.Close()

	body1, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if resp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update workspace"})
	}

	var updatedWorkspace models.TwWorkspace
	err = json.Unmarshal(body1, &updatedWorkspace)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	WorkspaceLog := models.TwWorkspaceLog{
		WorkspaceId:     workspaceUser.WorkspaceId,
		WorkspaceUserId: workspaceUser.ID,
		Action:          "Update",
		FieldChanged:    "Title, Description",
		OldValue:        Workspace.Title + ", " + Workspace.Description,
		NewValue:        updatedWorkspace.Title + ", " + updatedWorkspace.Description,
		Description:     "Update workspace",
	}
	_, err = dms.CallAPI(
		"POST",
		"/workspace_log",
		WorkspaceLog,
		nil,
		nil,
		120,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(updatedWorkspace)
}
