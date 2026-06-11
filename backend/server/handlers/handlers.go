// Package handlers holds the HTTP handlers. Each is a thin translator: decode →
// validate → call the app → translate the error → write JSON. Handlers share
// one Handlers type and the response helpers below.
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Simaky/go-github-tracker/backend/app"
)

// RequestIDKey is the gin-context key under which the RequestID middleware
// stores the per-request id.
const RequestIDKey = "request_id"

// Handlers holds the dependencies every handler shares.
type Handlers struct {
	app *app.App
}

// New constructs the Handlers.
func New(appInst *app.App) *Handlers {
	return &Handlers{app: appInst}
}

// Uptime is a plain-text liveness probe; it returns 200 once the process is up.
func (h *Handlers) Uptime(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

// RequestIDFrom returns the request id stored in the gin context, if any.
func RequestIDFrom(c *gin.Context) string {
	if v, ok := c.Get(RequestIDKey); ok {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}

// writeJSON writes v as a JSON response with the given status. Used by resource
// handlers as they are added.
func (h *Handlers) writeJSON(c *gin.Context, status int, v any) {
	c.JSON(status, v)
}

// writeError maps any error to the sanitised client envelope, logging the cause
// of 5xx responses with the request id.
func (h *Handlers) writeError(c *gin.Context, err error) {
	apiErr := asAPIError(err)
	if apiErr.Status >= http.StatusInternalServerError && apiErr.Cause != nil {
		log.Printf("request %s failed: code=%s cause=%s", RequestIDFrom(c), apiErr.Code, apiErr.Cause)
	}

	var body errorBody
	body.Error.Code = apiErr.Code
	body.Error.Message = apiErr.Message
	body.Error.Details = apiErr.Details
	c.JSON(apiErr.Status, body)
}

// decodeJSONBody decodes the request body into dst, rejecting unknown fields so
// client typos surface immediately.
func (h *Handlers) decodeJSONBody(c *gin.Context, dst any) error {
	dec := json.NewDecoder(c.Request.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return BadRequest(CodeInvalidRequest, "malformed request body")
	}
	return nil
}
