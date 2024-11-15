package notification_setting

import (
	"api/dms"
	"encoding/json"
	"errors"
	"github.com/timewise-team/timewise-models/models"
)

type NotificationSettingService struct {
}

func NewNotificationSettingService() *NotificationSettingService {

	return &NotificationSettingService{}
}

func (s NotificationSettingService) GetNotificationSettingByUserId(id string) (models.TwNotificationSettings, error) {
	var notificationSetting models.TwNotificationSettings

	resp, err := dms.CallAPI(
		"GET",
		"/notification_setting/"+id,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return notificationSetting, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&notificationSetting); err != nil {
		return notificationSetting, errors.New("error parsing response")
	}

	return notificationSetting, nil

}

func (s NotificationSettingService) UpdateNotificationSetting(id string, setting models.TwNotificationSettings) (models.TwNotificationSettings, error) {
	var notificationSetting models.TwNotificationSettings
	resp, err := dms.CallAPI(
		"PUT",
		"/notification_setting/"+id,
		setting,
		nil,
		nil,
		120,
	)
	if err != nil {
		return notificationSetting, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&notificationSetting); err != nil {
		return notificationSetting, err
	}

	return notificationSetting, nil

}
