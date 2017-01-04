package jwt

import (
	"fmt"
	"net/http"
	"github.com/codegangsta/negroni"
	"encoding/json"

	jwt "github.com/dgrijalva/jwt-go"
	request "github.com/dgrijalva/jwt-go/request"
  helperJwt "github.com/sfauvart/Agathadmin-api/helpers/jwt"
)

func RequireTokenAuthentication(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	authBackend := helperJwt.InitJWTAuthenticationBackend()

	token, err := request.ParseFromRequest(req, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		} else {
			return authBackend.PublicKey, nil
		}
	})

	if err == nil && token.Valid {
		next(rw, req)
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
	}
}

func RequireRoles(roles []string) negroni.HandlerFunc {
    return func(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
			authBackend := helperJwt.InitJWTAuthenticationBackend()

			token, err := request.ParseFromRequest(req, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				} else {
					return authBackend.PublicKey, nil
				}
			})

			if err == nil {
				var roleOk = false

				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					claimsJson, err := json.Marshal(&claims)
					if err != nil {
						fmt.Println(err)
					} else {
						var claimsCustom helperJwt.CustomClaims
						err = json.Unmarshal(claimsJson, &claimsCustom)

						for _, roleFromToken := range claimsCustom.User.Roles {
							for _, roleFromRoute := range roles {
								if roleFromToken == roleFromRoute {
									roleOk = true
								}
							}
						}
					}
				}

				if roleOk == true {
					next(rw, req)
				} else {
					rw.WriteHeader(http.StatusUnauthorized)
				}

			} else {
				rw.WriteHeader(http.StatusUnauthorized)
			}
		}
}
