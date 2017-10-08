package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/TerrenceHo/CalHacks4-Backend/models"
	"github.com/gorilla/mux"
)

func NewClasses(classes models.ClassService) *Classes {
	return &Classes{
		cs: classes,
	}
}

type Classes struct {
	cs models.ClassService
}

func (c *Classes) Create(w http.ResponseWriter, r *http.Request) {
	form := ClassesCreateForm{}
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	class := models.Class{
		Name:    form.Name,
		Summary: form.Summary,
		Videos:  form.Videos,
	}

	if err := c.cs.CreateClass(&class); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&class); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type ClassesCreateForm struct {
	Name    string   `json:"Name,omitempty"`
	Summary string   `json:"Summary,omitempty"`
	Videos  []string `json:"Videos,omitempty"`
}

func (c *Classes) GetAllClasses(w http.ResponseWriter, r *http.Request) {
	classes, err := c.cs.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(&classes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Classes) GetClass(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	class, err := c.cs.GetClass(vars["class"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(class); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
