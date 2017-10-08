package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/TerrenceHo/CalHacks4-Backend/models"
	"github.com/gorilla/mux"
)

func NewClasses(classes models.ClassService, videos models.VideoService) *Classes {
	return &Classes{
		cs: classes,
		vs: videos,
	}
}

type Classes struct {
	cs models.ClassService
	vs models.VideoService
}

func (c *Classes) Create(w http.ResponseWriter, r *http.Request) {
	form := ClassesCreateForm{}
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	class := models.Class{
		Name:        form.Name,
		Description: form.Description,
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
	Name        string `json:"Name,omitempty"`
	Description string `json:"Description,omitempty"`
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

// Used to upload videos into classes
func (c *Classes) Upload(w http.ResponseWriter, r *http.Request) {
	form := UploadForm{}
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uploads := form.Videos
	class_name := form.ClassName

	for i := 0; i < len(uploads); i++ {
		Video_URL := strings.Replace(uploads[i].Audio_URL, ".wav", ".mp4", 1)
		class, err := c.cs.GetClass(class_name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		video := models.Video{
			ClassID:           class.ID,
			URL:               Video_URL,
			Topics:            uploads[i].Topics,
			Related_Resources: uploads[i].Related_Resources,
		}
		err = c.vs.Create(&video)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := json.NewEncoder(w).Encode(&uploads); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type UploadForm struct {
	ClassName string            `json:"ClassName,omitempty"`
	Videos    []UploadAudioForm `json:Videos,omitempty"`
}

type UploadAudioForm struct {
	Audio_URL         string   `json:"Audio_URL,omitempty"`
	Topics            []string `json:"Topics,omitempty"`
	Related_Resources []string `json:"Related_Resources,omitempty"`
}

func (c *Classes) GetByKeyword(w http.ResponseWriter, r *http.Request) {
	keywords := []string{}
	if err := json.NewDecoder(r.Body).Decode(&keywords); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(&keywords); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
