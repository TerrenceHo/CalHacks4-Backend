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
		Password:      form.Email,
		PasswordReset: false,
	}

	if err := json.NewEncoder(w).Encode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// if err := u.us.Create(
}

type UsersCreateForm struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	UserType string `json:"usertype,omitempty"`
}
