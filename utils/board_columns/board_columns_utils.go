package board_columns

import (
	"errors"
	dtos "github.com/timewise-team/timewise-models/dtos/core_dtos/board_columns_dtos"
)

func ValidateBoardColumn(boardColumn dtos.BoardColumnsRequest) error {
	// Validate board column name
	if boardColumn.Name == "" {
		return errors.New("Board column name is required")
	}
	if len(boardColumn.Name) > 50 {
		return errors.New("Board column name must not exceed 50 characters")
	}
	if boardColumn.Position == 0 {
		return errors.New("Board column order is required")
	}
	return nil
}

func ValidateBoardColumnName(name string) error {
	// Validate board column name
	if name == "" {
		return errors.New("Board column name is required")
	}
	if len(name) > 50 {
		return errors.New("Board column name must not exceed 50 characters")
	}
	return nil
}
