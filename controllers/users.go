package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/TerrenceHo/CalHacks4-Backend/models"
)

func NewUsers(users models.UserService) *Users {
	return &Users{
		us: users,
	}
}

type Users struct {
	us models.UserService
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	form := UsersCreateForm{}
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := models.User{
		Name:          form.Name,
		UserType:      form.UserType,
		Email:         form.Email,
		Password:      form.Password,
		PasswordReset: false,
	}

	if err := u.us.Create(&user); err != nil {
		if pErr, ok := err.(PublicError); ok {
			http.Error(w, pErr.Public(), http.StatusNotAcceptable)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := json.NewEncoder(w).Encode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type UsersCreateForm struct {
	Name     string `json:"Name,omitempty"`
	Email    string `json:"Email,omitempty"`
	Password string `json:"Password,omitempty"`
	UserType string `json:"UserType,omitempty"`
}

func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	form := LoginForm{}
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		if pErr, ok := err.(PublicError); ok {
			http.Error(w, pErr.Public(), http.StatusNotAcceptable)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := json.NewEncoder(w).Encode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type LoginForm struct {
	Email    string `json:"Email,omitempty"`
	Password string `json:"Password,omitempty"`
}
