package auth

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	logger "github.com/Sirupsen/logrus"
	"github.com/sfauvart/Agathadmin-api/dao"
	helperJwt "github.com/sfauvart/Agathadmin-api/helpers/jwt"
	hJson "github.com/sfauvart/Agathadmin-api/helpers"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type SignInForm struct {
	Email    string `json:"email" form:"email"`
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
		hJson.SendJSONError(w, "Error while retrieving users", http.StatusInternalServerError)
		return
	}

	// Fetch user
	u, err := userDao.GetByEmail(requestSignIn.Email)
	if err != nil {
		hJson.SendJSONError(w, "Error while retrieving users", http.StatusNotFound)
		return
	}
	user := *u
	password := []byte(requestSignIn.Password)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), password)

	if err != nil {
		hJson.SendJSONError(w, "Error while retrieving users", http.StatusNotFound)
		return
	}

	authBackend := helperJwt.InitJWTAuthenticationBackend()
	token, err := authBackend.GenerateToken(user)
	if err != nil {
		hJson.SendJSONError(w, "Error while retrieving users", http.StatusNotFound)
		return
	} else {
		now := bson.Now()
		user.LastLogin = &now
		_, err = userDao.UpsertByID("", user)
		hJson.SendJSONOk(w, helperJwt.TokenAuthentication{token})
		return
	}
}
