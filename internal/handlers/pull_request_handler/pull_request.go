package pull_request_handler

import (
	"encoding/json"
	"errors"
	"github.com/SitnikovArtem06/avito-test-task/internal/model"
	"github.com/SitnikovArtem06/avito-test-task/internal/service/pull_request_service"
	"net/http"
)

type PullRequestHandler struct {
	prs *pull_request_service.PullRequestService
}

func NewPullRequestHAndler(prs *pull_request_service.PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{prs: prs}
}

func (p *PullRequestHandler) CreatePR(w http.ResponseWriter, r *http.Request) {

	var reqDto createPRDto

	json.NewDecoder(r.Body).Decode(&reqDto)

	if reqDto.PullRequestID == "" || reqDto.PullRequestName == "" || reqDto.AuthorID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]errorBodyDTO{
			"error": {
				Code:    "BAD_REQUEST",
				Message: "missing required fields",
			},
		})
		return
	}

	req := &model.CreatePullRequestRequest{
		Id:       reqDto.PullRequestID,
		Name:     reqDto.PullRequestName,
		AuthorId: reqDto.AuthorID,
	}

	pr, err := p.prs.CreatePullRequest(r.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, pull_request_service.ErrNotFoundTeam) || errors.Is(err, pull_request_service.ErrNotFoundAuthor):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]errorBodyDTO{
				"error": {
					Code:    "NOT_FOUND",
					Message: "resource not found",
				},
			})
		case errors.Is(err, pull_request_service.ErrAlreadyExist):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]errorBodyDTO{
				"error": {
					Code:    "PR_EXISTS",
					Message: "PR id already exists",
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

	resp := pullRequestDTO{
		PullRequestID:     pr.Id,
		PullRequestName:   pr.Name,
		AuthorID:          pr.AuthorId,
		Status:            pr.Status.String(),
		AssignedReviewers: nil,
	}

	for _, u := range pr.Reviewers {
		resp.AssignedReviewers = append(resp.AssignedReviewers, u.Id)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]pullRequestDTO{
		"pr": resp,
	})

}

func (p *PullRequestHandler) MergePR(w http.ResponseWriter, r *http.Request) {

	var req mergeRequest

	json.NewDecoder(r.Body).Decode(&req)

	if req.PullRequestId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]errorBodyDTO{
			"error": {
				Code:    "BAD_REQUEST",
				Message: "pull_request_id is required",
			},
		})
		return
	}

	pr, err := p.prs.Merge(r.Context(), req.PullRequestId)
	if err != nil {
		switch {
		case errors.Is(err, pull_request_service.ErrNotFound):
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

	resp := mergeResp{
		PullRequestID:     pr.Id,
		PullRequestName:   pr.Name,
		AuthorID:          pr.AuthorId,
		Status:            pr.Status.String(),
		AssignedReviewers: nil,
		MergedAt:          *pr.MergedAt,
	}

	for _, u := range pr.Reviewers {
		resp.AssignedReviewers = append(resp.AssignedReviewers, u.Id)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]mergeResp{
		"pr": resp,
	})

}

func (p *PullRequestHandler) ReassignReviewer(w http.ResponseWriter, r *http.Request) {

	var req reassignRequest

	json.NewDecoder(r.Body).Decode(&req)

	if req.PullRequestID == "" || req.OldReviewerID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]errorBodyDTO{
			"error": {
				Code:    "BAD_REQUEST",
				Message: "missing required fields",
			},
		})
		return
	}

	resp, err := p.prs.ReassignReviewer(r.Context(), req.PullRequestID, req.OldReviewerID)

	if err != nil {
		switch {
		case errors.Is(err, pull_request_service.ErrNotFoundUser) || errors.Is(err, pull_request_service.ErrNotFound) || errors.Is(err, pull_request_service.ErrNotFoundTeam):

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]errorBodyDTO{
				"error": {
					Code:    "NOT_FOUND",
					Message: "resource not found",
				},
			})
		case errors.Is(err, pull_request_service.ErrNotAssigned):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]errorBodyDTO{
				"error": {
					Code:    "NOT_ASSIGNED",
					Message: "old reviewer is not assigned"},
			})

		case errors.Is(err, pull_request_service.ErrMerged):

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]errorBodyDTO{
				"error": {
					Code:    "PR_MERGED",
					Message: "cannot reassign on merged PR"},
			})
		case errors.Is(err, pull_request_service.ErrNoCandidate):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]errorBodyDTO{
				"error": {
					Code:    "NO_CANDIDATE",
					Message: "no candidate found"},
			})

		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]errorBodyDTO{
				"error": {
					Code:    "INTERNAL",
					Message: "internal error",
				},
			})
		}
		return
	}

	prDto := pullRequestDTO{
		PullRequestID:     resp.Pr.Id,
		PullRequestName:   resp.Pr.Name,
		AuthorID:          resp.Pr.AuthorId,
		Status:            resp.Pr.Status.String(),
		AssignedReviewers: nil,
	}

	for _, u := range resp.Pr.Reviewers {
		prDto.AssignedReviewers = append(prDto.AssignedReviewers, u.Id)
	}

	respDto := reassignRespDto{
		Pr:          prDto,
		NewReviewer: resp.NewReviewer,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respDto)
}
