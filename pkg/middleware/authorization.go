package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/girishg4t/app_invite_service/pkg/model"
)

type contextKey string

var (
	AuthCtxKey contextKey = "auth-ctx-key"
)

// Authenticate check if the logged-in user is valid
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			http.Error(w, "token not present", http.StatusUnauthorized)
			return
		}

		var mySigningKey = []byte(os.Getenv("ACCESS_SECRET"))

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing")
			}
			return mySigningKey, nil
		})

		if err != nil {
			http.Error(w, "token has been expired", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["role"] == "ADMIN" {
				r.Header.Set("Role", "ADMIN")
				ctx := context.WithValue(context.Background(), AuthCtxKey, model.Token{
					Role:     "ADMIN",
					Username: fmt.Sprintf("%v", claims["username"]),
				})
				next.ServeHTTP(w, r.WithContext(ctx))
				return

			}
		}
		http.Error(w, "not authorized", http.StatusUnauthorized)
	})
}

// GetAuthorizationToken get's the token from the request context
func GetAuthorizationToken(ctx context.Context) (*model.Token, error) {
	val := ctx.Value(AuthCtxKey)
	if val == nil {
		return nil, errors.New("no permission")
	}
	u, ok := val.(model.Token)
	if !ok {
		return nil, errors.New("no permission")
	}
	return &u, nil
}
