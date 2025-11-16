package team_handler

import (
	"encoding/json"
	"errors"
	"github.com/SitnikovArtem06/avito-test-task/internal/service/team_service"
	"net/http"
)

type TeamHandler struct {
	ts *team_service.TeamService
}

func NewTeamHandler(ts *team_service.TeamService) *TeamHandler {
	return &TeamHandler{ts: ts}
}

func (t *TeamHandler) AddTeam(w http.ResponseWriter, r *http.Request) {

	var reqDto teamCreateRequestDTO

	json.NewDecoder(r.Body).Decode(&reqDto)

	if reqDto.TeamName == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]errorBodyDTO{
			"error": {
				Code:    "BAD_REQUEST",
				Message: "team_name is required",
			},
		})
		return
	}

	req := toCreateRequest(reqDto)

	team, err := t.ts.CreateTeam(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, team_service.ErrAlreadyExists):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]errorBodyDTO{
				"error": {
					Code:    "TEAM_EXISTS",
					Message: "team_name already exists",
				},
			})
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]errorBodyDTO{
				"error": {
					Code:    "INTERNAL",
					Message: "internal error",
				},
			})
		}

		return
	}

	resp := toRespDto(team)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]teamResponseDTO{
		"team": resp,
	})

}

func (t *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {

	teamName := r.URL.Query().Get("team_name")

	if teamName == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]errorBodyDTO{
			"error": {
				Code:    "BAD_REQUEST",
				Message: "team_name is required",
			},
		})
		return
	}

	team, err := t.ts.GetTeam(r.Context(), teamName)

	if err != nil {
		switch {
		case errors.Is(err, team_service.ErrNotFound):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]errorBodyDTO{
				"error": {
					Code:    "NOT_FOUND",
					Message: "resource not found",
				},
			})

		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]errorBodyDTO{
				"error": {
					Code:    "INTERNAL",
					Message: "internal error",
				},
			})
		}
		return
	}

	resp := toRespDto(team)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
