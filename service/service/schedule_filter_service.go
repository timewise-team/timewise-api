package service

import (
	"api/dms"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type ScheduleFilterService struct {
}

func NewScheduleFilterService() *ScheduleFilterService {
	return &ScheduleFilterService{}
}
func (s *ScheduleFilterService) ScheduleFilter(c *fiber.Ctx) (*http.Response, error) {
	param := c.Params("param")
	var resp *http.Response
	var err error
	if param == "" {
		resp, err = dms.CallAPI("GET", "/dbms/v1/schedule", nil, nil, nil, 120)
	} else {
		queryParams := map[string]string{"": param}
		resp, err = dms.CallAPI("GET", "/dbms/v1/schedule", nil, nil, queryParams, 120)
	}
	if err != nil {
		return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return resp, nil
}
