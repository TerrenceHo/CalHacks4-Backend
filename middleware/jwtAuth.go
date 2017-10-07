package middleware

import (
	"context"
	"net/http"

	"github.com/TerrenceHo/CalHacks4-Backend/config"
	"github.com/TerrenceHo/CalHacks4-Backend/controllers"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

type RequireJWT struct {
	VerifyKey []byte
	SignKey   []byte
}

func NewRequireJWT(conf *config.Config) *RequireJWT {
	return &RequireJWT{
		VerifyKey: conf.VerifyKey,
		SignKey:   conf.SignKey,
	}
}

func (rj *RequireJWT) AuthMW(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := controllers.Claims{}
		token, err := request.ParseFromRequestWithClaims(r, request.AuthorizationHeaderExtractor, &claims, func(token *jwt.Token) (interface{}, error) {
			verifyKeyRSA, err := jwt.ParseRSAPublicKeyFromPEM(rj.VerifyKey)
			if err != nil {
				panic(err)
			}
			return verifyKeyRSA, nil
		})
		if err != nil {
			switch err.(type) {
			case *jwt.ValidationError:
				vErr := err.(*jwt.ValidationError)
				switch vErr.Errors {
				case jwt.ValidationErrorExpired:
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				default:
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if token.Valid {
			newRequest := r.WithContext(context.WithValue(r.Context(), "user_claims", &claims))
			*r = *newRequest
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}
