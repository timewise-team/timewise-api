package transport

import (
	"api/service/notification_setting"
	"github.com/gofiber/fiber/v2"
)

type NotificationSettingHandler struct {
	service notification_setting.NotificationSettingService
}

func NewNotificationSettingHandler() *NotificationSettingHandler {
	service := notification_setting.NewNotificationSettingService()
	return &NotificationSettingHandler{
		service: *service,
	}
}

// getNotificationSettingByUserId godoc
// @Summary Get notification setting by user id
// @Description Get notification setting by user id
// @Tags notification_setting
// @Accept json
// @Produce json
// @Security bearerToken
// @Success 200 {object} models.TwNotificationSettings
// @Router /api/v1/notification_setting [get]
func (h NotificationSettingHandler) GetNotificationSettingByUserId(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userid").(string)
	notificationSetting, err := h.service.GetNotificationSettingByUserId(userId)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
	}
	return ctx.JSON(notificationSetting)
}

type UpdateNotificationSettingRequest struct {
	NotificationOnTag            bool
	NotificationOnComment        bool
	NotificationOnDueDate        bool
	NotificationOnScheduleChange bool
	NotificationOnEmail          bool
}

// UpdateNotificationSetting godoc
// @Summary Update notification setting
// @Description Update notification setting
// @Tags notification_setting
// @Accept json
// @Produce json
// @Param notification_setting body UpdateNotificationSettingRequest true "Notification Setting"
// @Success 200 {object} models.TwNotificationSettings
// @Security bearerToken
// @Router	/api/v1/notification_setting [put]
func (h NotificationSettingHandler) UpdateNotificationSetting(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userid").(string)
	var updateNotificationSettingRequest UpdateNotificationSettingRequest
	if err := ctx.BodyParser(&updateNotificationSettingRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	oldNotificationSetting, err := h.service.GetNotificationSettingByUserId(userId)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
	}
	oldNotificationSetting.NotificationOnTag = updateNotificationSettingRequest.NotificationOnTag
	oldNotificationSetting.NotificationOnComment = updateNotificationSettingRequest.NotificationOnComment
	oldNotificationSetting.NotificationOnDueDate = updateNotificationSettingRequest.NotificationOnDueDate
	oldNotificationSetting.NotificationOnScheduleChange = updateNotificationSettingRequest.NotificationOnScheduleChange
	oldNotificationSetting.NotificationOnEmail = updateNotificationSettingRequest.NotificationOnEmail

	notificationResponse, err := h.service.UpdateNotificationSetting(userId, oldNotificationSetting)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return ctx.JSON(notificationResponse)

}
