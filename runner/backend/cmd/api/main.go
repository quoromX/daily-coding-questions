package main

import (
	"fmt"
	"net/http"
	"os"

	authsvc "github.com/quoromx/irondoj/backend/internal/auth"
	"github.com/quoromx/irondoj/backend/internal/config"
	"github.com/quoromx/irondoj/backend/internal/database"
	"github.com/quoromx/irondoj/backend/internal/judge"
	"github.com/quoromx/irondoj/backend/internal/platform"
	"github.com/quoromx/irondoj/backend/internal/problems"
	"github.com/quoromx/irondoj/backend/internal/submissions"
	"github.com/quoromx/irondoj/backend/internal/users"
)

func main() {
	cfg := config.Load()
	store := database.NewStore()
	auth := authsvc.New(store, cfg.JWTAccessSecret, cfg.PasswordPepper)
	judgeClient := judge.New(cfg.Judge0BaseURL, cfg.Judge0AuthToken)

	userHandlers := users.New(store, auth)
	problemHandlers := problems.New(store)
	submissionHandlers := submissions.New(store, auth, judgeClient)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", platform.Wrap(func(w http.ResponseWriter, _ *http.Request) error {
		return platform.JSON(w, http.StatusOK, map[string]string{"status": "ok", "app": cfg.AppName})
	}))
	mux.HandleFunc("POST /v1/auth/register", platform.Wrap(userHandlers.Register))
	mux.HandleFunc("POST /v1/auth/login", platform.Wrap(userHandlers.Login))
	mux.HandleFunc("POST /v1/auth/refresh", platform.Wrap(userHandlers.Refresh))
	mux.HandleFunc("POST /v1/auth/logout", platform.Wrap(userHandlers.Logout))
	mux.HandleFunc("GET /v1/me", platform.Wrap(userHandlers.Me))
	mux.HandleFunc("GET /v1/me/dashboard", platform.Wrap(userHandlers.Dashboard))
	mux.HandleFunc("GET /v1/problems", platform.Wrap(problemHandlers.List))
	mux.HandleFunc("GET /v1/problems/", platform.Wrap(problemHandlers.Detail))
	mux.HandleFunc("POST /v1/problems/", platform.Wrap(func(w http.ResponseWriter, r *http.Request) error {
		if hasSuffix(r.URL.Path, "/run") {
			return submissionHandlers.Run(w, r)
		}
		if hasSuffix(r.URL.Path, "/submit") {
			return submissionHandlers.Submit(w, r)
		}
		return platform.Error{Status: http.StatusNotFound, Message: "route not found"}
	}))
	mux.HandleFunc("GET /v1/submissions", platform.Wrap(submissionHandlers.ListMine))
	mux.HandleFunc("GET /v1/submissions/", platform.Wrap(submissionHandlers.Detail))

	handler := platform.CORS(cfg.CORSOrigins)(platform.RequestLogger(mux))
	addr := ":" + cfg.Port
	fmt.Printf("%s API listening on %s\n", cfg.AppName, addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		fmt.Fprintf(os.Stderr, "server failed: %v\n", err)
		os.Exit(1)
	}
}

func hasSuffix(value, suffix string) bool {
	if len(value) < len(suffix) {
		return false
	}
	return value[len(value)-len(suffix):] == suffix
}
