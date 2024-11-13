package transport

import (
	"bytes"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/timewise-team/timewise-models/models"
	"net/http/httptest"
	"testing"
)

// Định nghĩa ScheduleService interface và MockScheduleService
type ScheduleService interface {
	GetScheduleById(id string) (*models.TwSchedule, error)
	CreateSchedule(schedule *models.TwSchedule) (*models.TwSchedule, error)
}

type MockScheduleService struct {
	mock.Mock
}

func (m *MockScheduleService) GetScheduleById(id string) (*models.TwSchedule, error) {
	args := m.Called(id)
	if schedule, ok := args.Get(0).(*models.TwSchedule); ok {
		return schedule, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockScheduleService) CreateSchedule(schedule *models.TwSchedule) (*models.TwSchedule, error) {
	args := m.Called(schedule)
	if createdSchedule, ok := args.Get(0).(*models.TwSchedule); ok {
		return createdSchedule, args.Error(1)
	}
	return nil, args.Error(1)
}

// Định nghĩa ScheduleHandler với ScheduleService
type ScheduleHandlerTest struct {
	service ScheduleService
}

func (h *ScheduleHandlerTest) GetScheduleByID(c *fiber.Ctx) error {
	scheduleID := c.Params("scheduleId")
	schedule, err := h.service.GetScheduleById(scheduleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(schedule)
}

func (h *ScheduleHandlerTest) CreateSchedule(c *fiber.Ctx) error {
	var schedule models.TwSchedule
	if err := c.BodyParser(&schedule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	createdSchedule, err := h.service.CreateSchedule(&schedule)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdSchedule)
}

// Hàm test thành công
func TestGetScheduleByID_Success(t *testing.T) {
	scheduleID := "1"
	mockSchedule := &models.TwSchedule{ID: 1, Title: "Meeting"}

	mockService := new(MockScheduleService)
	mockService.On("GetScheduleById", scheduleID).Return(mockSchedule, nil)

	app := fiber.New()
	handler := &ScheduleHandlerTest{service: mockService}
	app.Get("/schedules/:scheduleId", handler.GetScheduleByID)

	req := httptest.NewRequest("GET", "/schedules/1", nil)
	resp, err := app.Test(req)
	if assert.NoError(t, err) {
		defer resp.Body.Close()
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	}
	mockService.AssertExpectations(t)
}

// Hàm test khi xảy ra lỗi
func TestGetScheduleByID_Error(t *testing.T) {
	scheduleID := "999"
	mockError := errors.New("schedule not found")

	mockService := new(MockScheduleService)
	mockService.On("GetScheduleById", scheduleID).Return(nil, mockError)

	app := fiber.New()
	handler := &ScheduleHandlerTest{service: mockService}
	app.Get("/schedules/:scheduleId", handler.GetScheduleByID)

	req := httptest.NewRequest("GET", "/schedules/999", nil)
	resp, err := app.Test(req)
	if assert.NoError(t, err) {
		defer resp.Body.Close()
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	}
	mockService.AssertExpectations(t)
}

// Hàm test thành công cho CreateSchedule
func TestCreateSchedule_Success(t *testing.T) {
	mockService := new(MockScheduleService)

	mockService.On("CreateSchedule", mock.MatchedBy(func(schedule *models.TwSchedule) bool {
		return schedule.Title == "New Meeting" && schedule.Description == "Important meeting"
	})).Return(&models.TwSchedule{ID: 2, Title: "New Meeting", Description: "Important meeting"}, nil)

	app := fiber.New()
	handler := &ScheduleHandlerTest{service: mockService}
	app.Post("/schedules", handler.CreateSchedule)

	req := httptest.NewRequest("POST", "/schedules", bytes.NewReader([]byte(`{"title":"New Meeting","description":"Important meeting"}`)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if assert.NoError(t, err) {
		defer resp.Body.Close()
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	}
	mockService.AssertExpectations(t)
}

// Hàm test khi xảy ra lỗi cho CreateSchedule
func TestCreateSchedule_Error(t *testing.T) {
	mockService := new(MockScheduleService)
	mockError := errors.New("failed to create schedule")

	mockService.On("CreateSchedule", mock.MatchedBy(func(schedule *models.TwSchedule) bool {
		return schedule.Title == "New Meeting"
	})).Return(nil, mockError)

	app := fiber.New()
	handler := &ScheduleHandlerTest{service: mockService}
	app.Post("/schedules", handler.CreateSchedule)

	req := httptest.NewRequest("POST", "/schedules", bytes.NewReader([]byte(`{"title":"New Meeting"}`)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if assert.NoError(t, err) {
		defer resp.Body.Close()
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	}
	mockService.AssertExpectations(t)
}
