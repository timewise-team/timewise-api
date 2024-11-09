package transport

import (
	"api/service/board_columns"
	"api/service/comment"
	"api/service/document"
	"api/service/schedule"
	"api/service/schedule_participant"
	"github.com/gofiber/fiber/v2"
	schedule_dtos "github.com/timewise-team/timewise-models/dtos/core_dtos"
	dtos "github.com/timewise-team/timewise-models/dtos/core_dtos/board_columns_dtos"
	"github.com/timewise-team/timewise-models/models"
	"strconv"
)

// getBoardColumnsByWorkspace godoc
// @Summary Get all board columns by workspace (X-User-Email required, X-Workspace-Id required)
// @Description Get all board columns by workspace (X-User-Email required, X-Workspace-Id required)
// @Tags board_columns
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Param X-User-Email header string true "User email"
// @Param X-Workspace-Id header string true "Workspace ID"
// @Security bearerToken
// @Success 200 {array} []dtos.BoardColumnsResponse
// @Router /api/v1/board_columns/workspace/{workspace_id} [get]
func (h *BoardColumnsHandler) getBoardColumnsByWorkspace(c *fiber.Ctx) error {
	// Parse the request
	workspaceID := c.Params("workspace_id")
	if workspaceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workspace ID",
		})
	}
	workspaceIDCheck, err := strconv.Atoi(c.Params("workspace_id"))
	if err != nil || workspaceIDCheck <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workspace ID",
		})
	}
	workspaceUser := c.Locals("workspace_user").(*models.TwWorkspaceUser)
	if workspaceUser.WorkspaceId != workspaceIDCheck {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	// Get the board columns
	boardColumns, err := board_columns.NewBoardColumnsService().GetBoardColumnsByWorkspace(workspaceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if boardColumns == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "The server failed to respond",
		})
	}
	var boardColumnsResponseList []dtos.BoardColumnsResponse
	for _, boardColumn := range boardColumns {
		var boardColumnsResponse dtos.BoardColumnsResponse
		boardColumnsResponse.ID = boardColumn.ID
		boardColumnsResponse.Name = boardColumn.Name
		boardColumnsResponse.WorkspaceId = boardColumn.WorkspaceId
		boardColumnsResponse.Position = boardColumn.Position
		boardColumnsResponse.CreatedAt = boardColumn.CreatedAt
		boardColumnsResponse.UpdatedAt = boardColumn.UpdatedAt
		boardColumnsResponse.DeletedAt = boardColumn.DeletedAt

		schedules, err := schedule.NewScheduleService().GetSchedulesByBoardColumn(workspaceID, boardColumn.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "The server failed to respond",
			})
		}
		var schedulesListResponse []schedule_dtos.TwScheduleListInBoardColumnResponse
		if workspaceUser.Role == "admin" || workspaceUser.Role == "owner" {
			for _, schedule := range schedules {
				var schedulesList schedule_dtos.TwScheduleListInBoardColumnResponse
				schedulesList.ID = schedule.ID
				schedulesList.WorkspaceID = schedule.WorkspaceId
				schedulesList.BoardColumnID = schedule.BoardColumnId
				schedulesList.Title = schedule.Title
				schedulesList.Description = schedule.Description
				if schedule.StartTime != nil {
					schedulesList.StartTime = *schedule.StartTime
				}
				if schedule.EndTime != nil {
					schedulesList.EndTime = *schedule.EndTime
				}
				schedulesList.Location = schedule.Location
				schedulesList.CreatedBy = schedule.CreatedBy
				if schedule.CreatedAt != nil {
					schedulesList.CreatedAt = *schedule.CreatedAt
				}
				if schedule.UpdatedAt != nil {
					schedulesList.UpdatedAt = *schedule.UpdatedAt
				}
				schedulesList.Status = schedule.Status
				schedulesList.AllDay = schedule.AllDay
				schedulesList.Visibility = schedule.Visibility
				schedulesList.VideoTranscript = schedule.VideoTranscript
				schedulesList.ExtraData = schedule.ExtraData
				schedulesList.IsDeleted = schedule.IsDeleted
				schedulesList.RecurrencePattern = schedule.RecurrencePattern
				schedulesList.Position = schedule.Position
				schedulesList.Priority = schedule.Priority
				scheduleParticipants, err := schedule_participant.NewScheduleParticipantService().GetScheduleParticipantsBySchedule(schedule.ID, workspaceID)
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"message": "The server failed to respond",
					})
				}
				Documents, err := document.NewDocumentService().GetDocumentsBySchedule(schedule.ID)
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"message": "The server failed to respond",
					})
				}
				Comments, err := comment.NewCommentService().GetCommentsBySchedule(schedule.ID)
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"message": "The server failed to respond",
					})
				}
				schedulesList.ScheduleParticipants = scheduleParticipants
				schedulesList.Documents = len(Documents)
				schedulesList.Comments = len(Comments)

				schedulesListResponse = append(schedulesListResponse, schedulesList)

			}
		} else if workspaceUser.Role == "member" {
			for _, schedule := range schedules {
				if schedule.Visibility == "public" || schedule.Visibility == "" {
					var schedulesList schedule_dtos.TwScheduleListInBoardColumnResponse
					schedulesList.ID = schedule.ID
					schedulesList.WorkspaceID = schedule.WorkspaceId
					schedulesList.BoardColumnID = schedule.BoardColumnId
					schedulesList.Title = schedule.Title
					schedulesList.Description = schedule.Description
					if schedule.StartTime != nil {
						schedulesList.StartTime = *schedule.StartTime
					}
					if schedule.EndTime != nil {
						schedulesList.EndTime = *schedule.EndTime
					}
					schedulesList.Location = schedule.Location
					schedulesList.CreatedBy = schedule.CreatedBy
					if schedule.CreatedAt != nil {
						schedulesList.CreatedAt = *schedule.CreatedAt
					}
					if schedule.UpdatedAt != nil {
						schedulesList.UpdatedAt = *schedule.UpdatedAt
					}
					schedulesList.Status = schedule.Status
					schedulesList.AllDay = schedule.AllDay
					schedulesList.Visibility = schedule.Visibility
					schedulesList.VideoTranscript = schedule.VideoTranscript
					schedulesList.ExtraData = schedule.ExtraData
					schedulesList.IsDeleted = schedule.IsDeleted
					schedulesList.RecurrencePattern = schedule.RecurrencePattern
					schedulesList.Position = schedule.Position
					schedulesList.Priority = schedule.Priority
					scheduleParticipants, err := schedule_participant.NewScheduleParticipantService().GetScheduleParticipantsBySchedule(schedule.ID, workspaceID)
					if err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
							"message": "The server failed to respond",
						})
					}
					Documents, err := document.NewDocumentService().GetDocumentsBySchedule(schedule.ID)
					if err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
							"message": "The server failed to respond",
						})
					}
					Comments, err := comment.NewCommentService().GetCommentsBySchedule(schedule.ID)
					if err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
							"message": "The server failed to respond",
						})
					}
					schedulesList.ScheduleParticipants = scheduleParticipants
					schedulesList.Documents = len(Documents)
					schedulesList.Comments = len(Comments)

					schedulesListResponse = append(schedulesListResponse, schedulesList)
				} else if schedule.Visibility == "private" {
					//check member is participant or creator of schedule
					scheduleParticipants, err := schedule_participant.NewScheduleParticipantService().GetScheduleParticipantsBySchedule(schedule.ID, workspaceID)
					if err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
							"message": "The server failed to respond",
						})
					}
					var isParticipant bool
					for _, scheduleParticipant := range scheduleParticipants {
						if scheduleParticipant.WorkspaceUserId == workspaceUser.ID && scheduleParticipant.InvitationStatus == "joined" {
							isParticipant = true
							break
						}
					}
					if schedule.CreatedBy == workspaceUser.ID || isParticipant {
						var schedulesList schedule_dtos.TwScheduleListInBoardColumnResponse
						schedulesList.ID = schedule.ID
						schedulesList.WorkspaceID = schedule.WorkspaceId
						schedulesList.BoardColumnID = schedule.BoardColumnId
						schedulesList.Title = schedule.Title
						schedulesList.Description = schedule.Description
						if schedule.StartTime != nil {
							schedulesList.StartTime = *schedule.StartTime
						}
						if schedule.EndTime != nil {
							schedulesList.EndTime = *schedule.EndTime
						}
						schedulesList.Location = schedule.Location
						schedulesList.CreatedBy = schedule.CreatedBy
						if schedule.CreatedAt != nil {
							schedulesList.CreatedAt = *schedule.CreatedAt
						}
						if schedule.UpdatedAt != nil {
							schedulesList.UpdatedAt = *schedule.UpdatedAt
						}
						schedulesList.Status = schedule.Status
						schedulesList.AllDay = schedule.AllDay
						schedulesList.Visibility = schedule.Visibility
						schedulesList.VideoTranscript = schedule.VideoTranscript
						schedulesList.ExtraData = schedule.ExtraData
						schedulesList.IsDeleted = schedule.IsDeleted
						schedulesList.RecurrencePattern = schedule.RecurrencePattern
						schedulesList.Position = schedule.Position
						schedulesList.Priority = schedule.Priority
						scheduleParticipants, err := schedule_participant.NewScheduleParticipantService().GetScheduleParticipantsBySchedule(schedule.ID, workspaceID)
						if err != nil {
							return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
								"message": "The server failed to respond",
							})
						}
						Documents, err := document.NewDocumentService().GetDocumentsBySchedule(schedule.ID)
						if err != nil {
							return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
								"message": "The server failed to respond",
							})
						}
						Comments, err := comment.NewCommentService().GetCommentsBySchedule(schedule.ID)
						if err != nil {
							return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
								"message": "The server failed to respond",
							})
						}
						schedulesList.ScheduleParticipants = scheduleParticipants
						schedulesList.Documents = len(Documents)
						schedulesList.Comments = len(Comments)

						schedulesListResponse = append(schedulesListResponse, schedulesList)
					} else {
						var schedulesList schedule_dtos.TwScheduleListInBoardColumnResponse
						schedulesList.ID = schedule.ID
						schedulesList.WorkspaceID = schedule.WorkspaceId
						schedulesList.BoardColumnID = schedule.BoardColumnId
						schedulesList.Title = schedule.Title
						schedulesList.Description = schedule.Description
						schedulesList.Visibility = schedule.Visibility
						schedulesList.Position = schedule.Position
						schedulesList.Status = schedule.Status
						schedulesList.ExtraData = "IsLocked"
						schedulesListResponse = append(schedulesListResponse, schedulesList)
					}

				}
			}
		}
		boardColumnsResponse.Schedules = schedulesListResponse
		boardColumnsResponseList = append(boardColumnsResponseList, boardColumnsResponse)

	}
	// Return the response
	return c.JSON(boardColumnsResponseList)
}
