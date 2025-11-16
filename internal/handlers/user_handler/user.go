package user_handler

import (
	"encoding/json"
	"errors"
	"github.com/SitnikovArtem06/avito-test-task/internal/model"
	"github.com/SitnikovArtem06/avito-test-task/internal/service/user_service"
	"net/http"
)

type UserHandler struct {
	us *user_service.UserService
}

func NewUserHandler(us *user_service.UserService) *UserHandler {
	return &UserHandler{us: us}
}

func (u *UserHandler) SetActive(w http.ResponseWriter, r *http.Request) {

	var reqDto setActiveRequest

	json.NewDecoder(r.Body).Decode(&reqDto)

	if reqDto.Id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]errorBodyDTO{
			"error": {
				Code:    "BAD_REQUEST",
				Message: "user_id is required",
			},
		})
		return
	}

	req := &model.SetActiveRequest{
		Id:       reqDto.Id,
		IsActive: reqDto.IsActive,
	}

	user, err := u.us.UpdateActive(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, user_service.ErrNotFound):
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

	resp := userResp{
		Id:       user.Id,
		Name:     user.Name,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]userResp{
		"user": resp,
	})

}

func (u *UserHandler) GetPRsByUser(w http.ResponseWriter, r *http.Request) {

	userId := r.URL.Query().Get("user_id")

	if userId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]errorBodyDTO{
			"error": {
				Code:    "BAD_REQUEST",
				Message: "user_id is required",
			},
		})
		return
	}

	prs, err := u.us.GetPRsByUser(r.Context(), userId)
	if err != nil {
		switch {
		case errors.Is(err, user_service.ErrNotFound):
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

	var resp userReviewsResponse
	resp.UserID = userId

	for _, p := range prs {
		resp.PullRequests = append(resp.PullRequests, userReviewPullRequest{
			PullRequestID:   p.Id,
			PullRequestName: p.Name,
			AuthorID:        p.AuthorId,
			Status:          p.Status.String(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
