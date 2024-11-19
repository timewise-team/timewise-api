package notification

import (
	"api/dms"
	"github.com/timewise-team/timewise-models/models"
)

func PushNotifications(notification models.TwNotifications) error {
	// call dms to insert notification into database
	_, err := dms.CallAPI("POST", "/notification", notification, nil, nil, 120)
	if err != nil {
		return err
	}
	return nil
}
