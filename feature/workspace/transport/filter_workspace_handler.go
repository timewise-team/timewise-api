package transport

import (
	"api/dms"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"io/ioutil"
	"net/http"
)

// filterWorkspaces godoc
// @Summary Filter workspaces
// @Description Filter workspaces
// @Tags workspace
// @Accept json
// @Produce json
// @Security bearerToken
// @Param email query string false "Email"
// @Param role query string false "Role"
// @Param search query string false "Search"
// @Param sortBy query string false "Sort by"
// @Param order query string false "Order"
// @Success 200 {array} models.TwWorkspace
// @Router /api/v1/workspace/filter-workspaces [get]
func (h *WorkspaceHandler) filterWorkspaces(c *fiber.Ctx) error {
	queryParams := make(map[string]string)

	if c.Locals("userid") == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid userId"})
	}
	queryParams["userid"] = c.Locals("userid").(string)
	// Filter by email
	if email := c.Query("email"); email != "" {
		queryParams["email"] = email
	}

	// Filter by role
	if role := c.Query("role"); role != "" {
		queryParams["role"] = role
	}

	// Search by keyword
	if search := c.Query("search"); search != "" {
		queryParams["search"] = search
	}

	// Sort by field
	if sortBy := c.Query("sortBy"); sortBy != "" {
		order := c.Query("order", "asc")
		queryParams["sortBy"] = sortBy
		queryParams["order"] = order
	}

	// Call the API
	response, err := dms.CallAPI(
		"GET",
		"/workspace/filter/workspace",
		nil,
		nil,
		queryParams,
		120,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer response.Body.Close()

	// Read and parse the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if response.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get workspaces"})
	}

	var workspaces []models.TwWorkspace
	err = json.Unmarshal(body, &workspaces)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(workspaces)
}
