package middleware

import (
	"fmt"
	"net/http"

	"github.com/ultimateboy/discapi/utils"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/rs/rest-layer/resource"
)

// NewJWTHandler parse and validates JWT token if present and store it in the net/context
func NewJWTHandler(users *resource.Resource, jwtKeyFunc jwt.Keyfunc) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := request.ParseFromRequest(r, request.OAuth2Extractor, jwtKeyFunc)
			if err == request.ErrNoTokenInRequest {
				// If no token is found, let REST Layer hooks decide if the resource is public or not
				next.ServeHTTP(w, r)
				return
			}
			if err != nil || !token.Valid {
				if err != nil {
					fmt.Printf("error %s,", err)
				} else if !token.Valid {
					fmt.Printf("invalid token")
				}
				// @todo may want to return JSON error
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			claims := token.Claims.(jwt.MapClaims)
			userID, ok := claims["sub"].(string)
			if !ok || userID == "" {
				// The provided token is malformed, subject claim is missing
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			// Lookup the user by its id
			ctx := r.Context()
			user, err := users.Get(ctx, userID)
			if user != nil && err == resource.ErrUnauthorized {
				// Ignore unauthorized errors set by ourselves (see hooks.AuthUser)
				err = nil
			}
			if err != nil {
				// If user resource storage handler returned an error, respond with an error
				if err == resource.ErrNotFound {
					http.Error(w, "Invalid credentials", http.StatusForbidden)
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			// Store it into the request's context
			ctx = utils.NewContextWithUser(ctx, user)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
