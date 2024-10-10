package board_columns

import (
	"api/dms"
	"encoding/json"
	"errors"
	"fmt"
	"net/http" // Đảm bảo import gói http để sử dụng http.StatusOK và http.StatusCreated

	dtos "github.com/timewise-team/timewise-models/dtos/core_dtos/board_columns_dtos"
	"github.com/timewise-team/timewise-models/models"
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
		"/workspace/"+workspaceID+"/board_columns",
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
