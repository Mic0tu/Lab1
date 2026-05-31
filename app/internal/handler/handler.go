package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"notebook/internal/domain"
	"notebook/internal/repository"
)

type Handler struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("GET /api/users", h.listUsers)
	mux.HandleFunc("POST /api/users", h.createUser)
	mux.HandleFunc("GET /api/users/{id}", h.getUser)
	mux.HandleFunc("PUT /api/users/{id}", h.updateUser)
	mux.HandleFunc("DELETE /api/users/{id}", h.deleteUser)
	mux.HandleFunc("GET /api/profiles", h.listProfiles)
	mux.HandleFunc("POST /api/profiles", h.createProfile)
	mux.HandleFunc("GET /api/profiles/{id}", h.getProfile)
	mux.HandleFunc("PUT /api/profiles/{id}", h.updateProfile)
	mux.HandleFunc("DELETE /api/profiles/{id}", h.deleteProfile)
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.ListUsers(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list users")
		return
	}
	if users == nil {
		users = []domain.User{}
	}
	writeJSON(w, http.StatusOK, users)
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if input.Username == "" || input.Email == "" {
		writeError(w, http.StatusBadRequest, "username and email are required")
		return
	}

	u := domain.User{
		ID:       uuid.New(),
		Username: input.Username,
		Email:    input.Email,
	}
	if err := h.repo.CreateUser(r.Context(), u); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create user")
		return
	}
	writeJSON(w, http.StatusCreated, u)
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	u, err := h.repo.GetUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get user")
		return
	}
	writeJSON(w, http.StatusOK, u)
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if input.Username == "" || input.Email == "" {
		writeError(w, http.StatusBadRequest, "username and email are required")
		return
	}

	u := domain.User{ID: id, Username: input.Username, Email: input.Email}
	if err := h.repo.UpdateUser(r.Context(), u); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update user")
		return
	}
	writeJSON(w, http.StatusOK, u)
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	if err := h.repo.DeleteUser(r.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete user")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) listProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := h.repo.ListProfiles(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list profiles")
		return
	}
	if profiles == nil {
		profiles = []domain.Profile{}
	}
	writeJSON(w, http.StatusOK, profiles)
}

func (h *Handler) createProfile(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string `json:"name"`
		Lname string `json:"lname"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if input.Name == "" || input.Lname == "" {
		writeError(w, http.StatusBadRequest, "name and lname are required")
		return
	}

	p := domain.Profile{
		ID:    uuid.New(),
		Name:  input.Name,
		Lname: input.Lname,
	}
	if err := h.repo.CreateProfile(r.Context(), p); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create profile")
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid profile id")
		return
	}
	p, err := h.repo.GetProfile(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "profile not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get profile")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *Handler) updateProfile(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid profile id")
		return
	}
	var input struct {
		Name  string `json:"name"`
		Lname string `json:"lname"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if input.Name == "" || input.Lname == "" {
		writeError(w, http.StatusBadRequest, "name and lname are required")
		return
	}

	p := domain.Profile{ID: id, Name: input.Name, Lname: input.Lname}
	if err := h.repo.UpdateProfile(r.Context(), p); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "profile not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update profile")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *Handler) deleteProfile(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid profile id")
		return
	}
	if err := h.repo.DeleteProfile(r.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "profile not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete profile")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(strings.TrimSpace(s))
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
