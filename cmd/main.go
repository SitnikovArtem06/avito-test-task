package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/SitnikovArtem06/avito-test-task/internal/handlers"
	"github.com/SitnikovArtem06/avito-test-task/internal/handlers/pull_request_handler"
	"github.com/SitnikovArtem06/avito-test-task/internal/handlers/team_handler"
	"github.com/SitnikovArtem06/avito-test-task/internal/handlers/user_handler"
	"github.com/SitnikovArtem06/avito-test-task/internal/repository/pull_request_repository"
	"github.com/SitnikovArtem06/avito-test-task/internal/repository/team_repository"
	"github.com/SitnikovArtem06/avito-test-task/internal/repository/user_repository"
	"github.com/SitnikovArtem06/avito-test-task/internal/service/pull_request_service"
	"github.com/SitnikovArtem06/avito-test-task/internal/service/team_service"
	"github.com/SitnikovArtem06/avito-test-task/internal/service/user_service"
	"github.com/SitnikovArtem06/avito-test-task/pkg/database"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const TimeOut = 5

func run(ctx context.Context, port string) error {

	dbpool, err := database.InitDb(context.Background())

	if err != nil {
		return fmt.Errorf("Unable to connect to database: %v\n", err)
	}

	defer dbpool.Close()

	teamRep := team_repository.NewTeamRepository(dbpool)
	prRepo := pull_request_repository.NewPullRequestRepository(dbpool)
	userRep := user_repository.NewUserRepository(dbpool)

	teamService := team_service.NewTeamService(teamRep)
	userService := user_service.NewUserService(userRep, prRepo)
	prService := pull_request_service.NewPullRequestService(teamRep, userRep, prRepo)

	teamHandler := team_handler.NewTeamHandler(teamService)
	userHandler := user_handler.NewUserHandler(userService)
	prHandler := pull_request_handler.NewPullRequestHAndler(prService)

	r := handlers.Routes(teamHandler, userHandler, prHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	errCh := make(chan error, 1)

	go func() {
		log.Printf("Server start on : %s", port)
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():

		shutdownCtx, cancel := context.WithTimeout(context.Background(), TimeOut*time.Second)

		defer cancel()

		log.Println("Shutting down service-courier")

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown error: %w", err)
		}
		return nil

	case err = <-errCh:
		return fmt.Errorf("listen: %w", err)
	}
}
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found, using environment variables")
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer stop()

	if err := run(ctx, port); err != nil {
		log.Printf("fatal: %v", err)
		os.Exit(1)
	}
}
