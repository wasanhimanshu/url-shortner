package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.handleRegisterUser).Methods("POST")
	router.HandleFunc("/login", WithJWTAuth(h.handleLoginUser, h.db)).Methods("POST")
}

func (h *Handler) handleRegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("missing payload"))
		return
	}
	defer r.Body.Close()
	var payload *CreateUserPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}
	validate := validator.New()
	err = validate.Struct(payload)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		WriteError(w, http.StatusBadRequest, errors)
		return
	}
	var user = User{}
	rows := h.db.Where("email = ?", payload.Email).First(&user)
	if rows.RowsAffected > 0 {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("user already exist with email"))
		return
	} else {
		hashedPass, err := CreateHashPass(payload.Password)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, err)
			return
		}
		user = User{Firstname: payload.Firstname, Lastname: payload.Lastname, Email: payload.Email, Password: hashedPass}
		rows := h.db.Create(&user)
		if rows.Error != nil {
			WriteError(w, http.StatusInternalServerError, rows.Error)
		}
		token, err := CreateAuthandSetCookie(int(user.ID), w)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, err)
			return
		}
		WriteJSON(w, http.StatusCreated, token)
	}
}

func (h *Handler) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("missing payload"))
		return
	}
	defer r.Body.Close()
	var payload *LoginUserPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}
	validate := validator.New()
	err = validate.Struct(payload)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		WriteError(w, http.StatusBadRequest, errors)
		return
	}
	var user = User{}
	rows := h.db.Where("email = ?", payload.Email).First(&user)
	if rows.Error != nil {
		WriteError(w, http.StatusBadRequest, rows.Error)
		return
	}
	err = decryptHashPass(payload.Password, user.Password)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}
	token := r.Header.Get("Authorization")
	if token == "" {
		_, err := CreateAuthandSetCookie(int(user.ID), w)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, err)
			return
		}
	}
	WriteJSON(w, http.StatusAccepted, token)
}

func CreateAuthandSetCookie(id int, w http.ResponseWriter) (string, error) {
	trystring := "mypassword"
	secret := []byte(trystring)
	token, err := createJwt(secret, id)
	if err != nil {
		return "", err
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	return token, nil
}
