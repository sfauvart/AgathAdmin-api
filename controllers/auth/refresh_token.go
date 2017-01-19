package auth

import (
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	request "github.com/dgrijalva/jwt-go/request"
	"github.com/sfauvart/Agathadmin-api/helpers"
	helperJwt "github.com/sfauvart/Agathadmin-api/helpers/jwt"
	"net/http"
)

func RefreshTokenController(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	authBackend := helperJwt.InitJWTAuthenticationBackend()

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		} else {
			return authBackend.PublicKey, nil
		}
	})

	if err == nil && token.Valid {
		newToken, err := authBackend.RefreshToken(*token)
		if err != nil {
			helpers.SendJSONWithHTTPCode(w, helpers.Error{helpers.ErrorData{http.StatusUnauthorized, "refresh_token.errors.expired", "Your session is expired"}}, http.StatusUnauthorized)
			return
		} else {
			response, _ := json.Marshal(helperJwt.TokenAuthentication{newToken})
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			fmt.Fprintf(w, "%s", response)
			return
		}
	} else {
		helpers.SendJSONWithHTTPCode(w, helpers.Error{helpers.ErrorData{http.StatusUnauthorized, "refresh_token.errors.expired", "Your session is expired"}}, http.StatusUnauthorized)
		return
	}
}
