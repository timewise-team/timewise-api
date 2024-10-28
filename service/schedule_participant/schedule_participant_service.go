package schedule_participant

import (
	"api/config"
	"api/dms"
	"api/service/auth"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/schedule_participant_dtos"
	"github.com/timewise-team/timewise-models/models"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ScheduleParticipantService struct {
}

func NewScheduleParticipantService() *ScheduleParticipantService {
	return &ScheduleParticipantService{}
}

func (h *ScheduleParticipantService) GetScheduleParticipantsBySchedule(scheduleId int, workspaceId string) ([]schedule_participant_dtos.ScheduleParticipantInfo, error) {
	scheduleIdStr := strconv.Itoa(scheduleId)
	if scheduleIdStr == "" {
		return nil, nil
	}
	if workspaceId == "" {
		return nil, nil
	}
	resp, err := dms.CallAPI(
		"GET",
		"/schedule_participant/workspace/"+workspaceId+"/schedule/"+scheduleIdStr,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var scheduleParticipants []schedule_participant_dtos.ScheduleParticipantInfo
	if err := json.NewDecoder(resp.Body).Decode(&scheduleParticipants); err != nil {
		return nil, err
	}

	return scheduleParticipants, nil

}

func (h *ScheduleParticipantService) GetScheduleParticipantsByScheduleID(scheduleId int) ([]schedule_participant_dtos.ScheduleParticipantInfo, error) {
	scheduleIdStr := strconv.Itoa(scheduleId)
	if scheduleIdStr == "" {
		return nil, nil
	}
	resp, err := dms.CallAPI(
		"GET",
		"/schedule_participant/schedule/"+scheduleIdStr,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var scheduleParticipants []schedule_participant_dtos.ScheduleParticipantInfo
	if err := json.NewDecoder(resp.Body).Decode(&scheduleParticipants); err != nil {
		return nil, err
	}

	return scheduleParticipants, nil

}

func (h *ScheduleParticipantService) InviteToSchedule(c *fiber.Ctx, InviteToScheduleDto schedule_participant_dtos.InviteToScheduleRequest) (*schedule_participant_dtos.ScheduleParticipantResponse, error) {

	// Lay ra workspaceUser cua thg gui
	workspaceUserInvite, ok := c.Locals("workspace_user").(*models.TwWorkspaceUser)
	if !ok {
		return nil, errors.New("Failed to retrieve schedule participant")
	}

	// Lay ra workspaceUser cua thg dc gui
	resp, err := dms.CallAPI(
		"GET",
		"/workspace_user/email/"+InviteToScheduleDto.Email+"/workspace/"+strconv.Itoa(workspaceUserInvite.WorkspaceId),
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var workspaceUserInvited models.TwWorkspaceUser
	if err := json.NewDecoder(resp.Body).Decode(&workspaceUserInvited); err != nil {
		return nil, err
	}

	// Lay ra scheduleParticipant cua thg gui
	scheduleParticipant, ok := c.Locals("scheduleParticipant").(models.TwScheduleParticipant)
	if !ok {
		return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve schedule participant",
		})
	}

	resp1, err := dms.CallAPI(
		"GET",
		"/schedule/"+strconv.Itoa(scheduleParticipant.ScheduleId),
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, fmt.Errorf("server error: %v", err)
	}
	defer resp1.Body.Close()

	if resp1.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp1.StatusCode)
	}

	// Lay schedule
	var schedule models.TwSchedule
	if err := json.NewDecoder(resp1.Body).Decode(&schedule); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	cfg, err1 := config.LoadConfig()
	if err1 != nil {
		return nil, c.Status(500).JSON(fiber.Map{
			"message": "Failed to load config",
		})
	}
	acceptLink, declineLink, _ := auth.GenerateInviteScheduleLinks(cfg, scheduleParticipant.ScheduleId, workspaceUserInvited.ID)

	var scheduleParticipantResponse schedule_participant_dtos.ScheduleParticipantResponse
	resp3, err := dms.CallAPI(
		"GET",
		"/schedule_participant/workspace_user/"+strconv.Itoa(workspaceUserInvited.ID)+
			"/schedule/"+strconv.Itoa(scheduleParticipant.ScheduleId),
		nil,
		nil,
		nil,
		120,
	)

	now := time.Now()

	if resp3.StatusCode == http.StatusNotFound {
		newScheduleParticipant := models.TwScheduleParticipant{
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			ScheduleId:       scheduleParticipant.ScheduleId,
			WorkspaceUserId:  workspaceUserInvited.ID,
			AssignBy:         workspaceUserInvite.ID,
			InvitationSentAt: &now,
			InvitationStatus: "pending",
		}

		resp2, err2 := dms.CallAPI(
			"POST",
			"/schedule_participant",
			newScheduleParticipant,
			nil,
			nil,
			120,
		)
		if err2 != nil {
			return nil, err2
		}
		defer resp2.Body.Close()

		body, err2 := ioutil.ReadAll(resp2.Body)
		if err2 != nil {
			return nil, errors.New("cannot read response body")
		}

		if resp.StatusCode != http.StatusOK {
			return nil, errors.New(string(body))
		}

		if errParsing := json.Unmarshal(body, &scheduleParticipantResponse); errParsing != nil {
			return nil, errors.New("error parsing JSON")
		}

		subject := fmt.Sprintf("Invitation to join schedule: %s", schedule.Title)
		content := auth.BuildScheduleInvitationContent(&schedule, acceptLink, declineLink)
		if err := auth.SendInvitationEmail(cfg, InviteToScheduleDto.Email, content, subject); err != nil {
			return nil, c.Status(500).JSON(fiber.Map{"message": "Failed to send invitation email"})
		}
	} else {
		if err != nil {
			log.Printf("Error calling API: %v", err)
			return nil, err
		}

		// Kiểm tra mã trạng thái phản hồi
		if resp3.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API returned status code %d", resp3.StatusCode)
		}

		// Đọc nội dung phản hồi
		body, err := io.ReadAll(resp3.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			return nil, err
		}
		defer resp3.Body.Close()

		// Giải mã JSON thành struct TwScheduleParticipant
		var scheduleParticipantInvited models.TwScheduleParticipant

		if err := json.Unmarshal(body, &scheduleParticipantInvited); err != nil {
			log.Printf("Error unmarshalling response: %v", err)
			return nil, err
		}

		switch scheduleParticipantInvited.InvitationStatus {
		case "joined":
			return nil, errors.New("User is already in the schedule")

		case "pending":
			subject := fmt.Sprintf("Reminder: Invitation to join schedule: %s", schedule.Title)
			content := auth.BuildScheduleInvitationContent(&schedule, acceptLink, declineLink)
			if err := auth.SendInvitationEmail(cfg, InviteToScheduleDto.Email, content, subject); err != nil {
				return nil, c.Status(500).JSON(fiber.Map{"message": "Failed to send invitation email"})
			}

		case "declined", "removed":
			subject := fmt.Sprintf("Reminder: Invitation to join schedule: %s", schedule.Title)
			content := auth.BuildScheduleInvitationContent(&schedule, acceptLink, declineLink)
			if err := auth.SendInvitationEmail(cfg, InviteToScheduleDto.Email, content, subject); err != nil {
				return nil, c.Status(500).JSON(fiber.Map{"message": "Failed to send invitation email"})
			}
			scheduleParticipantInvited.UpdatedAt = time.Now()
			scheduleParticipantInvited.InvitationStatus = "pending"

			resp4, err4 := dms.CallAPI(
				"PUT",
				"/schedule_participant/"+strconv.Itoa(scheduleParticipantInvited.ID),
				scheduleParticipantInvited,
				nil,
				nil,
				120,
			)

			if err4 != nil {
				return nil, err4
			}
			defer resp4.Body.Close()

			body, err2 := ioutil.ReadAll(resp4.Body)
			if err2 != nil {
				return nil, errors.New("cannot read response body")
			}

			if resp4.StatusCode != http.StatusOK {
				return nil, errors.New(string(body))
			}

			if errParsing := json.Unmarshal(body, &scheduleParticipantResponse); errParsing != nil {
				return nil, errors.New("error parsing JSON")
			}
		}
	}

	return &scheduleParticipantResponse, nil
}

func (h *ScheduleParticipantService) AcceptInvite(scheduleId, workspaceUserId string) (*schedule_participant_dtos.ScheduleParticipantResponse, error) {

	resp, err := dms.CallAPI(
		"GET",
		"/schedule_participant/workspace_user/"+workspaceUserId+"/schedule/"+scheduleId,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var scheduleParticipants models.TwScheduleParticipant
	if err := json.NewDecoder(resp.Body).Decode(&scheduleParticipants); err != nil {
		return nil, err
	}

	now := time.Now()
	scheduleParticipants.Status = "participant"
	scheduleParticipants.ResponseTime = &now
	scheduleParticipants.InvitationStatus = "joined"

	resp1, err := dms.CallAPI(
		"PUT",
		"/schedule_participant/"+strconv.Itoa(scheduleParticipants.ID),
		scheduleParticipants,
		nil,
		nil,
		120,
	)

	if err != nil {
		return nil, err
	}
	defer resp1.Body.Close()

	var updateScheduleParticipants schedule_participant_dtos.ScheduleParticipantResponse
	if err := json.NewDecoder(resp1.Body).Decode(&updateScheduleParticipants); err != nil {
		return nil, err
	}

	return &updateScheduleParticipants, nil

}

func (h *ScheduleParticipantService) DeclineInvite(scheduleId, workspaceUserId string) (*schedule_participant_dtos.ScheduleParticipantResponse, error) {

	resp, err := dms.CallAPI(
		"GET",
		"/schedule_participant/workspace_user/"+workspaceUserId+"/schedule/"+scheduleId,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var scheduleParticipants models.TwScheduleParticipant
	if err := json.NewDecoder(resp.Body).Decode(&scheduleParticipants); err != nil {
		return nil, err
	}

	now := time.Now()
	scheduleParticipants.ResponseTime = &now
	scheduleParticipants.InvitationStatus = "declined"

	resp1, err := dms.CallAPI(
		"PUT",
		"/schedule_participant/"+strconv.Itoa(scheduleParticipants.ID),
		scheduleParticipants,
		nil,
		nil,
		120,
	)

	if err != nil {
		return nil, err
	}
	defer resp1.Body.Close()

	var updateScheduleParticipants schedule_participant_dtos.ScheduleParticipantResponse
	if err := json.NewDecoder(resp1.Body).Decode(&updateScheduleParticipants); err != nil {
		return nil, err
	}

	return &updateScheduleParticipants, nil

}
