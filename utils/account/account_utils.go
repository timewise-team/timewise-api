package account

import "github.com/timewise-team/timewise-models/dtos/core_dtos"

func IsValidInputUpdateProfileRequest(request core_dtos.UpdateProfileRequestDto) bool {
	if request.FirstName == "" {
		return false
	}
	if request.LastName == "" {
		return false
	}
	return true
}
