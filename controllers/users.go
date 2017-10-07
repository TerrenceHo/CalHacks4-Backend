package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/TerrenceHo/CalHacks4-Backend/models"
	jwt "github.com/dgrijalva/jwt-go"
)

func NewUsers(users models.UserService, signKey []byte) *Users {
	return &Users{
		us:      users,
		signKey: signKey,
	}
}

type Users struct {
	us      models.UserService
	signKey []byte
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

	tokenString, err := u.createUserJWT(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&Token{tokenString}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type LoginForm struct {
	Email    string `json:"Email,omitempty"`
	Password string `json:"Password,omitempty"`
}

// form to send back JWT token
type Token struct {
	Token string `json:"token"`
}

// Takes user information and creates a JWT token with it, and signs the token
// with signKey.  Errors should never occur here, but if they do, then our app
// is in a really bad state.  Returns JWT token and nil
func (u *Users) createUserJWT(user *models.User) (string, error) {
	// Create claims for the jwt
	claims := Claims{
		user.Email,
		user.ID,
		jwt.StandardClaims{
			// ExpiresAt: time.Now().Add(time.Minute * 20).Unix(),
			Issuer: "user",
		},
	}

	//Sign the jwt
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signKeyRSA, err := jwt.ParseRSAPrivateKeyFromPEM(u.signKey)
	if err != nil {
		return "", err
	}
	tokenString, err := t.SignedString(signKeyRSA)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Struct to Handle jwt claims
type Claims struct {
	UserEmail string `json:"user_email,omitempty"`
	UserID    uint   `json:"user_id,omitempty"`
	jwt.StandardClaims
}

// Used on app open to check if a user is valid.  Middleware jwt should take
// care of everything, as jwt must be checked before this runs.
// Also sends back an array of user materials the user can send with the app.
func (u *Users) Check(w http.ResponseWriter, r *http.Request) {
	// claims := r.Context().Value("user_claims").(*Claims)
	// user, err := u.us.ByID(claims.UserID)
	// if err != nil {
	// 	if pErr, ok := err.(PublicError); ok {
	// 		http.Error(w, pErr.Public(), http.StatusNotFound)
	// 		return
	// 	} else {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// }
	// resp, err := json.Marshal(user.MaterialTypes)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }
	w.Header().Set("Content-Type", "application/json")
	// w.Write(resp)
}
