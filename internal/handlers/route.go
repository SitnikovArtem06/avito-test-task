package handlers

import (
	"github.com/SitnikovArtem06/avito-test-task/internal/handlers/pull_request_handler"
	"github.com/SitnikovArtem06/avito-test-task/internal/handlers/team_handler"
	"github.com/SitnikovArtem06/avito-test-task/internal/handlers/user_handler"
)
import "github.com/go-chi/chi/v5"

func Routes(th *team_handler.TeamHandler, uh *user_handler.UserHandler, ph *pull_request_handler.PullRequestHandler) chi.Router {

	r := chi.NewRouter()
	r.Post("/team/add", th.AddTeam)
	r.Get("/team/get", th.GetTeam)
	r.Post("/users/setIsActive", uh.SetActive)
	r.Post("/pullRequest/create", ph.CreatePR)
	r.Post("/pullRequest/merge", ph.MergePR)
	r.Post("/pullRequest/reassign", ph.ReassignReviewer)
	r.Get("/users/getReview", uh.GetPRsByUser)
	return r
}
