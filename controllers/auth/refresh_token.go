package auth

import (
	"net/http"
	"fmt"
	"encoding/json"
  jwt "github.com/dgrijalva/jwt-go"
	request "github.com/dgrijalva/jwt-go/request"
  helperJwt "github.com/sfauvart/Agathadmin-api/helpers/jwt"
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			response, _ := json.Marshal(helperJwt.TokenAuthentication{newToken})
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			fmt.Fprintf(w, "%s", response)
			return
		}
  } else {
    w.WriteHeader(http.StatusUnauthorized)
  }
}
