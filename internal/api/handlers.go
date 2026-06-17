package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/pratikkarelia25/gosnip/internal/shortcode"
	"github.com/pratikkarelia25/gosnip/internal/store"
	"github.com/pratikkarelia25/gosnip/internal/validate"
)

const (
	generatedCodeLength = 7
	maxGenerateAttempts = 5
)

type Handler struct {
	store *store.Store
}

func NewHandler(store *store.Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", h.root)
	r.Get("/{code}", h.redirect)
	r.Post("/shorten", h.shorten)
	return r
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func (h *Handler) root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to GoSnip!"))
}

func (h *Handler) redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	longUrl, err := h.store.GetActiveLongUrl(code)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "This snip doesnt exists!")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, longUrl, http.StatusTemporaryRedirect)
}

type shortenRequest struct {
	ShortCode        string `json:"short_code"`
	LongUrl          string `json:"long_url"`
	ExpiresInSeconds int64  `json:"expires_in_seconds"`
}

type shortenResponse struct {
	Message   string     `json:"message"`
	ShortCode string     `json:"short_code"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

func (h *Handler) shorten(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.LongUrl == "" {
		writeError(w, http.StatusBadRequest, "long_url is required")
		return
	}
	if !validate.URL(req.LongUrl) {
		writeError(w, http.StatusBadRequest, "wrong URL format")
		return
	}
	if req.ExpiresInSeconds < 0 {
		writeError(w, http.StatusBadRequest, "expires_in_seconds must not be negative")
		return
	}

	if existingCode, err := h.store.FindActiveShortCode(req.LongUrl); err == nil {
		writeJSON(w, http.StatusOK, shortenResponse{
			Message:   "Short URL already exists for this long URL",
			ShortCode: existingCode,
		})
		return
	} else if !errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusInternalServerError, "Failed to check for an existing short URL")
		return
	}

	code, err := h.resolveShortCode(req.ShortCode)
	if err != nil {
		if errors.Is(err, errCodeTaken) {
			writeError(w, http.StatusConflict, "short_code is already in use")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to generate a short code")
		return
	}

	var expiresAt *time.Time
	if req.ExpiresInSeconds > 0 {
		t := time.Now().Add(time.Duration(req.ExpiresInSeconds) * time.Second)
		expiresAt = &t
	}

	if err := h.store.InsertMapping(code, req.LongUrl, expiresAt); err != nil {
		fmt.Println("Error inserting mapping into database:", err)
		writeError(w, http.StatusInternalServerError, "Failed to create short URL")
		return
	}

	writeJSON(w, http.StatusCreated, shortenResponse{
		Message:   "Short URL created successfully",
		ShortCode: code,
		ExpiresAt: expiresAt,
	})
}

var errCodeTaken = errors.New("short code already in use")

// resolveShortCode returns requestedCode if it's free, generates a
// random unused one if requestedCode is empty, or errCodeTaken if
// requestedCode is already in use.
func (h *Handler) resolveShortCode(requestedCode string) (string, error) {
	if requestedCode != "" {
		taken, err := h.store.ShortCodeTaken(requestedCode)
		if err != nil {
			return "", err
		}
		if taken {
			return "", errCodeTaken
		}
		return requestedCode, nil
	}

	for i := 0; i < maxGenerateAttempts; i++ {
		code, err := shortcode.Generate(generatedCodeLength)
		if err != nil {
			return "", err
		}
		taken, err := h.store.ShortCodeTaken(code)
		if err != nil {
			return "", err
		}
		if !taken {
			return code, nil
		}
	}

	return "", fmt.Errorf("could not generate a unique short code after %d attempts", maxGenerateAttempts)
}
