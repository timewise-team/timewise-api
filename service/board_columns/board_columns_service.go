package board_columns

import (
	"api/dms"
	"api/service/workspace"
	board_columns_utils "api/utils/board_columns"
	"encoding/json"
	"fmt"
	dtos "github.com/timewise-team/timewise-models/dtos/core_dtos/board_columns_dtos"
	"github.com/timewise-team/timewise-models/models"
	"net/http" // Đảm bảo import gói http để sử dụng http.StatusOK và http.StatusCreated
	"strconv"
)

type BoardColumnsService struct {
}

func NewBoardColumnsService() *BoardColumnsService {
	return &BoardColumnsService{}
}

func (s *BoardColumnsService) CreateBoardColumn(request dtos.BoardColumnsRequest) (*models.TwBoardColumn, error) {
	if err := board_columns_utils.ValidateBoardColumn(request); err != nil {
		return nil, fmt.Errorf("invalid request: %v", err)
	}
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
	// Validate input
	if workspaceID == "" {
		return nil, fmt.Errorf("workspace id is required")
	}
	if _, err := strconv.Atoi(workspaceID); err != nil {
		return nil, fmt.Errorf("invalid workspace id")
	}
	result := workspace.NewWorkspaceService().GetWorkspaceById(workspaceID)
	if result == nil {
		return nil, fmt.Errorf("workspace not found")
	}
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

	return boardColumns, nil
}

func (h *BoardColumnsService) InitBoardColumns(workspaceID int) error {
	// workspace_service
	workspaceIDStr := fmt.Sprintf("%d", workspaceID)
	workspaceCheck := workspace.NewWorkspaceService().GetWorkspaceById(workspaceIDStr)
	if workspaceCheck == nil {
		return fmt.Errorf("workspace not found")
	}

	var boardColumns = []models.TwBoardColumn{
		{
			Name:        "Title",
			Position:    1,
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
	if boardColumnId == "" {
		return fmt.Errorf("board column id is required")
	}
	if _, err := strconv.Atoi(boardColumnId); err != nil {
		return fmt.Errorf("invalid board column id")
	}
	_, err := h.GetBoardColumnById(boardColumnId)
	if err != nil {
		return fmt.Errorf("failed to get board column: %v", err)
	}

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

func (h *BoardColumnsService) UpdateBoardColumn(boardColumnId string, request string) (*models.TwBoardColumn, error) {
	// Validate input
	if boardColumnId == "" {
		return nil, fmt.Errorf("board column id is required")
	}
	if request == "" {
		return nil, fmt.Errorf("request is required")
	}
	if len(request) > 50 {
		return nil, fmt.Errorf("request is too long")
	}
	_, err := h.GetBoardColumnById(boardColumnId)
	if err != nil {
		return nil, fmt.Errorf("board column not found: %v", err)
	}

	// Call API
	boardColumnRequest := models.TwBoardColumn{
		Name: request,
	}
	resp, err := dms.CallAPI(
		"PUT",
		"/board_columns/"+boardColumnId,
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

func (h *BoardColumnsService) GetBoardColumnById(boardColumnId string) (*models.TwBoardColumn, error) {
	if boardColumnId == "" {
		return nil, fmt.Errorf("board column id is required")
	}
	if _, err := strconv.Atoi(boardColumnId); err != nil {
		return nil, fmt.Errorf("invalid board column id")
	}
	// Call API
	resp, err := dms.CallAPI(
		"GET",
		"/board_columns/"+boardColumnId,
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

	var boardColumn models.TwBoardColumn
	if err := json.NewDecoder(resp.Body).Decode(&boardColumn); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &boardColumn, nil
}

func (s *BoardColumnsService) UpdatePositionsAfterDeletion(position int, id int) error {
	if position < 0 {
		return fmt.Errorf("invalid position")
	}
	if id < 0 {
		return fmt.Errorf("invalid workspace id")
	}
	if result := workspace.NewWorkspaceService().GetWorkspaceById(fmt.Sprintf("%d", id)); result == nil {
		return fmt.Errorf("workspace not found")
	}

	// Tạo payload cho API
	payload := map[string]interface{}{
		"position":     position,
		"workspace_id": id,
	}

	// Gọi API với phương thức PUT hoặc PATCH (tùy thuộc vào API của bạn)
	resp, err := dms.CallAPI(
		"PUT", // hoặc "PATCH", tùy vào API của bạn
		"/board_columns/update_position_after_deletion/position",
		payload,
		nil,
		nil,
		120,
	)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	defer resp.Body.Close()

	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *BoardColumnsService) UpdatePositionAfterDrag(oldPosition int, newPosition int, workspaceId int, board_column_id string) error {

	if newPosition < 0 {
		return fmt.Errorf("invalid position")
	}
	if oldPosition < 0 {
		return fmt.Errorf("invalid position")
	}
	if workspaceId < 0 {
		return fmt.Errorf("invalid workspace id")
	}
	if board_column_id == "" {
		return fmt.Errorf("board column id is required")
	}
	// Kiểm tra xem vị trí mới có khác với vị trí cũ không
	if oldPosition == newPosition {
		return nil
	}
	// Tạo một danh sách để chứa các cột cần cập nhật
	var columnsToUpdate []models.TwBoardColumn

	// Trường hợp kéo sang phải (newPosition > oldPosition)
	if newPosition > oldPosition {
		// Lấy các cột có position từ oldPosition + 1 đến newPosition và giảm position của chúng đi 1
		columns, err := s.GetColumnsInRange(oldPosition+1, newPosition, workspaceId)
		if err != nil {
			return fmt.Errorf("failed to get columns in range: %v", err)
		}

		// Giảm position của từng cột trong khoảng
		for _, column := range columns {
			column.Position -= 1
			columnsToUpdate = append(columnsToUpdate, column)
		}
	} else {
		// Trường hợp kéo sang trái (newPosition < oldPosition)
		// Lấy các cột có position từ newPosition đến oldPosition - 1 và tăng position của chúng lên 1
		columns, err := s.GetColumnsInRange(newPosition, oldPosition-1, workspaceId)
		if err != nil {
			return fmt.Errorf("failed to get columns in range: %v", err)
		}

		// Tăng position của từng cột trong khoảng
		for _, column := range columns {
			column.Position += 1
			columnsToUpdate = append(columnsToUpdate, column)
		}
	}

	// Cập nhật vị trí của cột đang di chuyển (đặt nó vào vị trí mới)
	draggedColumn, err := s.GetBoardColumnById(board_column_id)
	if err != nil {
		return fmt.Errorf("failed to get dragged column: %v", err)
	}
	draggedColumn.Position = newPosition
	columnsToUpdate = append(columnsToUpdate, *draggedColumn)

	// Gửi các cập nhật đến API hoặc cơ sở dữ liệu
	for _, column := range columnsToUpdate {
		err := s.UpdateBoardColumnPosition(column)
		if err != nil {
			return fmt.Errorf("failed to update column position: %v", err)
		}
	}

	return nil
}

func (s *BoardColumnsService) GetColumnsInRange(position1 int, position2 int, workspaceId int) ([]models.TwBoardColumn, error) {
	if position1 < 0 || position2 < 0 {
		return nil, fmt.Errorf("invalid position")
	}
	if workspaceId < 0 {
		return nil, fmt.Errorf("invalid workspace id")
	}
	if result := workspace.NewWorkspaceService().GetWorkspaceById(fmt.Sprintf("%d", workspaceId)); result == nil {
		return nil, fmt.Errorf("workspace not found")
	}
	// Tạo payload cho API
	payload := map[string]interface{}{
		"position1":    position1,
		"position2":    position2,
		"workspace_id": workspaceId,
	}

	// Gọi API với phương thức GET
	resp, err := dms.CallAPI(
		"GET",
		"/board_columns/range/position",
		payload,
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
	var columns []models.TwBoardColumn
	if err := json.NewDecoder(resp.Body).Decode(&columns); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return columns, nil
}

func (s *BoardColumnsService) UpdateBoardColumnPosition(column models.TwBoardColumn) error {
	//Validate input
	if column.ID < 0 {
		return fmt.Errorf("invalid column id")
	}
	if column.Position < 0 {
		return fmt.Errorf("invalid position")
	}
	if column.WorkspaceId < 0 {
		return fmt.Errorf("invalid workspace id")
	}
	if column.Name == "" {
		return fmt.Errorf("invalid name")
	}

	resp, err := dms.CallAPI(
		"PUT",
		"/board_columns/update_position/position",
		column,
		nil,
		nil,
		120,
	)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
