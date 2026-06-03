package submissions

import (
	"net/http"
	"strings"

	"github.com/quoromx/irondoj/backend/internal/auth"
	"github.com/quoromx/irondoj/backend/internal/database"
	"github.com/quoromx/irondoj/backend/internal/judge"
	"github.com/quoromx/irondoj/backend/internal/platform"
)

type Handlers struct {
	store *database.Store
	auth  *auth.Service
	judge *judge.Client
}

func New(store *database.Store, authService *auth.Service, judgeClient *judge.Client) *Handlers {
	return &Handlers{store: store, auth: authService, judge: judgeClient}
}

type submitRequest struct {
	LanguageID int    `json:"languageId"`
	Language   string `json:"language"`
	SourceCode string `json:"sourceCode"`
}

func (h *Handlers) Run(w http.ResponseWriter, r *http.Request) error {
	return h.evaluate(w, r, false)
}

func (h *Handlers) Submit(w http.ResponseWriter, r *http.Request) error {
	return h.evaluate(w, r, true)
}

func (h *Handlers) ListMine(w http.ResponseWriter, r *http.Request) error {
	user, ok := h.currentUser(r)
	if !ok {
		return platform.Error{Status: http.StatusUnauthorized, Message: "authentication required"}
	}
	return platform.JSON(w, http.StatusOK, map[string]any{"submissions": h.store.UserSubmissions(user.ID)})
}

func (h *Handlers) Detail(w http.ResponseWriter, r *http.Request) error {
	id := strings.TrimPrefix(r.URL.Path, "/v1/submissions/")
	submission, ok := h.store.Submission(id)
	if !ok {
		return platform.Error{Status: http.StatusNotFound, Message: "submission not found"}
	}
	return platform.JSON(w, http.StatusOK, map[string]any{"submission": submission})
}

func (h *Handlers) evaluate(w http.ResponseWriter, r *http.Request, hidden bool) error {
	user, ok := h.currentUser(r)
	if !ok {
		return platform.Error{Status: http.StatusUnauthorized, Message: "authentication required"}
	}
	slug := strings.TrimPrefix(r.URL.Path, "/v1/problems/")
	slug = strings.TrimSuffix(slug, "/run")
	slug = strings.TrimSuffix(slug, "/submit")
	problem, ok := h.store.ProblemWithTests(slug)
	if !ok {
		return platform.Error{Status: http.StatusNotFound, Message: "problem not found"}
	}
	var req submitRequest
	if err := platform.Decode(r, &req); err != nil {
		return platform.Error{Status: http.StatusBadRequest, Message: "invalid submission payload"}
	}
	if req.Language == "" {
		req.Language = "Go"
	}
	if req.LanguageID == 0 {
		req.LanguageID = 60
	}

	results := make([]database.SubmissionResult, 0)
	passed := 0
	for _, test := range problem.TestCases {
		if test.Hidden && !hidden {
			continue
		}
		result, err := h.judge.Run(r.Context(), req.SourceCode, req.LanguageID, test)
		if err != nil {
			return platform.Error{Status: http.StatusBadGateway, Message: "judge execution failed"}
		}
		results = append(results, result)
		if result.Status == "Passed" {
			passed++
		}
	}
	status := "Failed"
	if passed == len(results) && len(results) > 0 {
		status = "Passed"
	}
	submission := h.store.SaveSubmission(database.Submission{
		UserID:      user.ID,
		ProblemSlug: problem.Slug,
		ProblemName: problem.Title,
		LanguageID:  req.LanguageID,
		Language:    req.Language,
		SourceCode:  req.SourceCode,
		Status:      status,
		Score:       score(passed, len(results)),
		RuntimeMS:   passed * 12,
		MemoryKB:    2048,
		Results:     results,
	})
	return platform.JSON(w, http.StatusCreated, map[string]any{"submission": submission})
}

func (h *Handlers) currentUser(r *http.Request) (database.User, bool) {
	return h.auth.UserFromToken(platform.BearerToken(r))
}

func score(passed, total int) int {
	if total == 0 {
		return 0
	}
	return passed * 100 / total
}
