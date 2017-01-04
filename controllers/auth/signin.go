package auth

import (
	"net/http"
	"github.com/sfauvart/Agathadmin-api/dao"
	helperJwt "github.com/sfauvart/Agathadmin-api/helpers/jwt"
	"golang.org/x/crypto/bcrypt"
	"encoding/json"
	"fmt"
	logger "github.com/Sirupsen/logrus"
)

type SignInForm struct {
	Email string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type SignInResponse struct {
	Message string `json:"message"`
	IdToken string `json:"id_token"`
}

func SignInController(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	requestSignIn := new(SignInForm)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestSignIn)

	userDao, err := dao.GetUserDAO(dao.DAOMongo)
	if err != nil {
		logger.WithField("error", err).Fatal("unable to connect to mongo db")
		w.WriteHeader(500)
		return
	}

	// Fetch user
	u, err := userDao.GetByEmail(requestSignIn.Email)
	if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(404)
			return
	}
	user := *u
	password := []byte(requestSignIn.Password)
	tstPwd := bcrypt.CompareHashAndPassword([]byte(user.Password), password)

	authBackend := helperJwt.InitJWTAuthenticationBackend()
	token, err := authBackend.GenerateToken(user)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		return
	} else {
		response, _ := json.Marshal(helperJwt.TokenAuthentication{token})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, "%s", response)
		return
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
	fmt.Println(tstPwd)
}
