package notfication

import (
	"api/dms"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"io"
)

type NotificationService struct {
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (s *NotificationService) GetNotifications(userEmailIds []string) ([]models.TwNotifications, error) {
	resp, err := dms.CallAPI("POST", "/notification/user-email-ids", userEmailIds, nil, nil, 120)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != fiber.StatusOK {
		return nil, errors.New("error from external service: " + string(body))
	}

	var notifications []models.TwNotifications
	if err := json.Unmarshal(body, &notifications); err != nil {
		return nil, err
	}

	return notifications, nil
}
