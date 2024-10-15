package board_columns

import (
	"api/dms"
	"encoding/json"
	"errors"
	"fmt"
	dtos "github.com/timewise-team/timewise-models/dtos/core_dtos/board_columns_dtos"
	"github.com/timewise-team/timewise-models/models"
	"net/http" // Đảm bảo import gói http để sử dụng http.StatusOK và http.StatusCreated
)

type BoardColumnsService struct {
}

func NewBoardColumnsService() *BoardColumnsService {
	return &BoardColumnsService{}
}

func (s *BoardColumnsService) CreateBoardColumn(request dtos.BoardColumnsRequest) (*models.TwBoardColumn, error) {
	// Call API
	boardColumnRequest := models.TwBoardColumn{
		Name:        request.Name,
		Position:    request.Position,
		WorkspaceId: request.WorkspaceId,
	}
	resp, err := dms.CallAPI(
		"POST",
		"/board_columns",
		boardColumnRequest,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, fmt.Errorf("server error: %v", err)
	}
	defer resp.Body.Close()

	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var boardColumn models.TwBoardColumn
	if err := json.NewDecoder(resp.Body).Decode(&boardColumn); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &boardColumn, nil
}

func (s *BoardColumnsService) GetBoardColumnsByWorkspace(workspaceID string) ([]models.TwBoardColumn, error) {
	// Call API
	resp, err := dms.CallAPI(
		"GET",
		"/board_columns/workspace/"+workspaceID,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, fmt.Errorf("server error: %v", err)
	}
	defer resp.Body.Close()
	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var boardColumns []models.TwBoardColumn
	if err := json.NewDecoder(resp.Body).Decode(&boardColumns); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if len(boardColumns) == 0 {
		return nil, errors.New("no board columns found")
	}

	return boardColumns, nil
}

func (h *BoardColumnsService) InitBoardColumns(workspaceID int) error {
	var boardColumns = []models.TwBoardColumn{
		{
			Name:        "To Do",
			Position:    1,
			WorkspaceId: workspaceID,
		},
		{
			Name:        "In Progress",
			Position:    2,
			WorkspaceId: workspaceID,
		},
		{
			Name:        "Done",
			Position:    3,
			WorkspaceId: workspaceID,
		},
	}
	for _, boardColumn := range boardColumns {
		// Call API
		resp, err := dms.CallAPI(
			"POST",
			"/board_columns",
			boardColumn,
			nil,
			nil,
			120,
		)
		if err != nil {
			return fmt.Errorf("server error: %v", err)
		}
		defer resp.Body.Close()

		// Kiểm tra mã trạng thái HTTP
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

	}

	return nil
}

func (h *BoardColumnsService) DeleteBoardColumn(boardColumnId string) error {
	// Call API
	resp, err := dms.CallAPI(
		"DELETE",
		"/board_columns/"+string(boardColumnId),
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	defer resp.Body.Close()

	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (h *BoardColumnsService) UpdateBoardColumn(boardColumnId string, request dtos.BoardColumnsRequest) (*models.TwBoardColumn, error) {
	// Call API
	boardColumnRequest := models.TwBoardColumn{
		Name:     request.Name,
		Position: request.Position,
	}
	resp, err := dms.CallAPI(
		"PUT",
		"/board_columns/"+string(boardColumnId),
		boardColumnRequest,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, fmt.Errorf("server error: %v", err)
	}
	defer resp.Body.Close()

	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var boardColumn models.TwBoardColumn
	if err := json.NewDecoder(resp.Body).Decode(&boardColumn); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &boardColumn, nil
}
