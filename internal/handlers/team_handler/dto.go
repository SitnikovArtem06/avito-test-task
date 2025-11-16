package team_handler

import "github.com/SitnikovArtem06/avito-test-task/internal/model"

type teamMemberDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type teamCreateRequestDTO struct {
	TeamName string          `json:"team_name"`
	Members  []teamMemberDTO `json:"members"`
}

type teamResponseDTO struct {
	TeamName string          `json:"team_name"`
	Members  []teamMemberDTO `json:"members"`
}

type errorBodyDTO struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func toCreateRequest(reqDto teamCreateRequestDTO) *model.TeamCreateRequest {
	req := &model.TeamCreateRequest{
		TeamName: reqDto.TeamName,
		Members:  make([]model.User, len(reqDto.Members)),
	}
	for i, m := range reqDto.Members {
		req.Members[i] = model.User{
			Id:       m.UserID,
			Name:     m.Username,
			IsActive: m.IsActive,
		}
	}
	return req
}

func toRespDto(team *model.Team) teamResponseDTO {
	resp := teamResponseDTO{
		TeamName: team.Name,
		Members:  make([]teamMemberDTO, len(team.Members)),
	}
	for i, m := range team.Members {
		resp.Members[i] = teamMemberDTO{
			UserID:   m.Id,
			Username: m.Name,
			IsActive: m.IsActive,
		}
	}
	return resp
}
