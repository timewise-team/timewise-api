package transport

import (
	"api/middleware"
	"github.com/gofiber/fiber/v2"
)

type ScheduleParticipantHandlerRegister struct {
	Router  fiber.Router
	Handler *ScheduleParticipantHandler
}

func RegisterScheduleParticipantHandler(router fiber.Router) {
	handler := NewScheduleParticipantHandler()
	scheduleParticipantHandler := &ScheduleParticipantHandlerRegister{
		Handler: handler,
	}

	// Register all endpoints here
	router.Get("schedule/:scheduleId",
		scheduleParticipantHandler.Handler.GetScheduleParticipantByScheduleID)
	router.Post("invite",
		middleware.CheckWorkspaceRole([]string{"owner", "admin", "member"}),
		middleware.CheckScheduleStatus([]string{"creator"}),
		scheduleParticipantHandler.Handler.InviteToSchedule)
	router.Get("/accept-invitation-via-email/token/:token?", scheduleParticipantHandler.Handler.AcceptInvite)
	router.Get("/decline-invitation-via-email/token/:token?", scheduleParticipantHandler.Handler.DeclineInvite)
	router.Put("/assign",
		middleware.CheckWorkspaceRole([]string{"owner", "admin", "member"}),
		middleware.CheckScheduleStatus([]string{"creator"}),
		scheduleParticipantHandler.Handler.AssignMember)
	router.Put("/remove/:id",
		middleware.CheckWorkspaceRole([]string{"owner", "admin", "member"}),
		middleware.CheckScheduleStatus([]string{"creator"}),
		scheduleParticipantHandler.Handler.RemoveParticipant)
	router.Put("/unassign/:id",
		middleware.CheckWorkspaceRole([]string{"owner", "admin", "member"}),
		middleware.CheckScheduleStatus([]string{"creator"}),
		scheduleParticipantHandler.Handler.UnassignParticipant)
}
