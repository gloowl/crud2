package server

import (
	"crud2/internal/models"
	"crud2/internal/repository"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{
		repo: repository.NewUserRepository(db),
	}
}

func (u *UserHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users/register", u.RegisterUser).Methods("POST")
	r.HandleFunc("/users/login", u.LoginUser).Methods("POST")
}

// POST /users/register
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка парсинга JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	user := &models.User{
		UserName: req.UserName,
		Password: req.Password,
	}

	// Валидация
	if err := user.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка, не существует ли пользователь уже
	existing, err := h.repo.GetUserByUsername(user.UserName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при проверке пользователя: %v", err), http.StatusInternalServerError)
		return
	}
	if existing != nil {
		http.Error(w, "Пользователь уже существует", http.StatusConflict)
		return
	}

	// Хеш пароля
	if err := user.SetPassword(req.Password); err != nil {
		http.Error(w, "Ошибка хеширования пароля: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Создание
	if err := h.repo.CreateUser(user); err != nil {
		http.Error(w, "Ошибка при создании пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message":  "Пользователь успешно зарегистрирован",
		"username": user.UserName,
	})
}

// POST /users/login
func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка парсинга JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUserByUsername(req.UserName)
	if err != nil {
		http.Error(w, "Ошибка при получении пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	if !user.CheckPassword(req.Password) {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{
		"message":  "Успешный вход",
		"username": user.UserName,
	})
}
