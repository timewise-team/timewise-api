package transport

import (
	"api/service/account"
	"api/service/notfication"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type NotificationHandler struct {
	notificationService notfication.NotificationService
	accountService      account.AccountService
}

func NewNotificationHandler() *NotificationHandler {
	notificationService := notfication.NewNotificationService()
	accountService := account.NewAccountService()
	return &NotificationHandler{
		notificationService: *notificationService,
		accountService:      *accountService,
	}
}

// GetNotifications godoc
// @Summary Get notifications
// @Description Get notifications
// @Tags notification
// @Accept json
// @Produce json
// @Security bearerToken
// @Success 200 {array} models.TwNotifications
// @Router /api/v1/notification [get]
func (h NotificationHandler) GetNotifications(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userid").(string)

	userEmails, err := h.accountService.GetLinkedUserEmails(userId, "")
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).SendString("Error when GetLinkedUserEmails")
	}

	var userEmailIds []int
	for _, email := range userEmails {
		userEmailIds = append(userEmailIds, email.ID)
	}

	var userEmailIdsStr []string
	for _, id := range userEmailIds {
		userEmailIdsStr = append(userEmailIdsStr, strconv.Itoa(id))
	}

	notifications, err := h.notificationService.GetNotifications(userEmailIdsStr)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
	}
	return ctx.JSON(notifications)
}

// UpdateNotification godoc
// @Summary Update notification
// @Description Update notification
// @Tags notification
// @Accept json
// @Produce json
// @Security bearerToken
// @Param notification_id path string true "Notification ID"
// @Success 200 {string} string
// @Router /api/v1/notification/{notification_id} [put]
func (h NotificationHandler) UpdateNotification(ctx *fiber.Ctx) error {
	notificationId := ctx.Params("notification_id")

	err := h.notificationService.UpdateNotification(notificationId)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
	}
	return ctx.JSON("Update notification successfully")
}
