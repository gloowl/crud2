package server

import (
	"crud2/internal/models"
	"crud2/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type HerbHandler struct {
	repo *repository.HerbRepository
}

func NewHerbHandler(repo *repository.HerbRepository) *HerbHandler {
	return &HerbHandler{repo: repo}
}

// RegisterRoutes adds herb routes to router
func (h *HerbHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/herbs", h.GetAll).Methods("GET")
	r.HandleFunc("/herbs/{id:[0-9]+}", h.GetByID).Methods("GET")
	r.HandleFunc("/herbs", h.Create).Methods("POST")
	r.HandleFunc("/herbs/{id:[0-9]+}", h.Update).Methods("PUT")
	r.HandleFunc("/herbs/{id:[0-9]+}", h.Delete).Methods("DELETE")
	r.HandleFunc("/herbs/search", h.Search).Methods("GET")
	r.HandleFunc("/herbs/poisonous", h.GetPoisonous).Methods("GET")
}

// GET /herbs
func (h *HerbHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	herbs, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, herbs, http.StatusOK)
}

// GET /herbs/{id}
func (h *HerbHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "некорректный ID", http.StatusBadRequest)
		return
	}

	herb, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	jsonResponse(w, herb, http.StatusOK)
}

// POST /herbs
func (h *HerbHandler) Create(w http.ResponseWriter, r *http.Request) {
	var herb models.Herb

	if err := json.NewDecoder(r.Body).Decode(&herb); err != nil {
		http.Error(w, "ошибка парсинга JSON", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(&herb); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse(w, herb, http.StatusCreated)
}

// PUT /herbs/{id}
func (h *HerbHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "некорректный ID", http.StatusBadRequest)
		return
	}

	var herb models.Herb

	if err := json.NewDecoder(r.Body).Decode(&herb); err != nil {
		http.Error(w, "ошибка парсинга JSON", http.StatusBadRequest)
		return
	}

	herb.ID = id

	if err := h.repo.Update(&herb); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse(w, herb, http.StatusOK)
}

// DELETE /herbs/{id}
func (h *HerbHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "некорректный ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /herbs/search?name=...
func (h *HerbHandler) Search(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "параметр name обязателен", http.StatusBadRequest)
		return
	}

	herbs, err := h.repo.Search(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, herbs, http.StatusOK)
}

// GET /herbs/poisonous
func (h *HerbHandler) GetPoisonous(w http.ResponseWriter, r *http.Request) {
	herbs, err := h.repo.GetPoisonous()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, herbs, http.StatusOK)

}

// --- helpers ---
func parseID(idStr string) (int, error) {
	return strconv.Atoi(idStr)
}

func jsonResponse(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
