package notification

import (
	"api/dms"
	"errors"
	"github.com/timewise-team/timewise-models/models"
)

func PushNotifications(notification models.TwNotifications) error {
	// Validate required fields
	if notification.UserEmailId == 0 || notification.Type == "" || notification.Message == "" {
		return errors.New("Missing required fields")
	}

	// call dms to insert notification into database
	_, err := dms.CallAPI("POST", "/notification", notification, nil, nil, 120)
	if err != nil {
		return err
	}
	return nil
}
