package problems

import (
	"net/http"
	"strings"

	"github.com/quoromx/irondoj/backend/internal/database"
	"github.com/quoromx/irondoj/backend/internal/platform"
)

type Handlers struct {
	store *database.Store
}

func New(store *database.Store) *Handlers {
	return &Handlers{store: store}
}

func (h *Handlers) List(w http.ResponseWriter, r *http.Request) error {
	problems := h.store.Problems(r.URL.Query().Get("q"), r.URL.Query().Get("difficulty"), r.URL.Query().Get("tag"))
	return platform.JSON(w, http.StatusOK, map[string]any{"problems": problems})
}

func (h *Handlers) Detail(w http.ResponseWriter, r *http.Request) error {
	slug := strings.TrimPrefix(r.URL.Path, "/v1/problems/")
	if strings.Contains(slug, "/") {
		slug = strings.Split(slug, "/")[0]
	}
	problem, ok := h.store.Problem(slug)
	if !ok {
		return platform.Error{Status: http.StatusNotFound, Message: "problem not found"}
	}
	return platform.JSON(w, http.StatusOK, map[string]any{"problem": problem})
}
