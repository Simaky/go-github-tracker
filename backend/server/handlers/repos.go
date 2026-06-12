package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Simaky/go-github-tracker/backend/app"
)

// CreateRepo handles POST /api/repos: fetch from GitHub, persist, return 201.
func (h *Handlers) CreateRepo(c *gin.Context) {
	var req app.CreateRepoRequest
	if err := h.decodeJSONBody(c, &req); err != nil {
		h.writeError(c, err)
		return
	}
	repo, err := h.app.TrackRepo(c.Request.Context(), req)
	if err != nil {
		h.writeError(c, err)
		return
	}
	h.writeJSON(c, http.StatusCreated, repo)
}

// TotalMetrics handles GET /api/metrics.
func (h *Handlers) TotalMetrics(c *gin.Context) {
	metrics, err := h.app.TotalMetrics(c.Request.Context())
	if err != nil {
		h.writeError(c, err)
		return
	}
	h.writeJSON(c, http.StatusOK, metrics)
}

// ListRepos handles GET /api/repos, optionally filtered by ?language=.
func (h *Handlers) ListRepos(c *gin.Context) {
	repos, err := h.app.ListRepos(c.Request.Context(), c.Query("language"))
	if err != nil {
		h.writeError(c, err)
		return
	}
	h.writeJSON(c, http.StatusOK, repos)
}

// GetRepo handles GET /api/repos/:id.
func (h *Handlers) GetRepo(c *gin.Context) {
	id, err := repoID(c)
	if err != nil {
		h.writeError(c, err)
		return
	}
	repo, err := h.app.GetRepo(c.Request.Context(), id)
	if err != nil {
		h.writeError(c, err)
		return
	}
	h.writeJSON(c, http.StatusOK, repo)
}

// UpdateNotes handles PATCH /api/repos/:id.
func (h *Handlers) UpdateNotes(c *gin.Context) {
	id, err := repoID(c)
	if err != nil {
		h.writeError(c, err)
		return
	}
	var req app.UpdateNotesRequest
	if err := h.decodeJSONBody(c, &req); err != nil {
		h.writeError(c, err)
		return
	}
	repo, err := h.app.UpdateNotes(c.Request.Context(), id, req)
	if err != nil {
		h.writeError(c, err)
		return
	}
	h.writeJSON(c, http.StatusOK, repo)
}

// RefreshRepo handles POST /api/repos/:id/refresh.
func (h *Handlers) RefreshRepo(c *gin.Context) {
	id, err := repoID(c)
	if err != nil {
		h.writeError(c, err)
		return
	}
	repo, err := h.app.RefreshRepo(c.Request.Context(), id)
	if err != nil {
		h.writeError(c, err)
		return
	}
	h.writeJSON(c, http.StatusOK, repo)
}

// DeleteRepo handles DELETE /api/repos/:id.
func (h *Handlers) DeleteRepo(c *gin.Context) {
	id, err := repoID(c)
	if err != nil {
		h.writeError(c, err)
		return
	}
	if err := h.app.DeleteRepo(c.Request.Context(), id); err != nil {
		h.writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// repoID parses and validates the :id path parameter.
func repoID(c *gin.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return 0, BadRequest(CodeInvalidRequest, "id must be a positive integer")
	}
	return id, nil
}
