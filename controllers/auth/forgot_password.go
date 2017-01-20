package auth

import (
	"crypto/rand"
	"encoding/base64"
	logger "github.com/Sirupsen/logrus"
	"github.com/sfauvart/Agathadmin-api/dao"
	"github.com/sfauvart/Agathadmin-api/helpers"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type ForgotPasswordForm struct {
	Email string `json:"email"`
}
type CheckTokenForm struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Password string `json:"password"`
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func ForgotPasswordController(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var f ForgotPasswordForm
	err := helpers.GetJSONContent(&f, r)
	if err != nil {
		logger.WithField("error", err).WithField("form", f).Warn("ForgotPasswordController error")
		return
	}

	userDao, err := dao.GetUserDAO(dao.DAOMongo)
	if err != nil {
		logger.WithField("error", err).Fatal("unable to connect to mongo db")
		helpers.SendJSONError(w, "Error while sending forgot password token", http.StatusInternalServerError)
		return
	}

	user, err := userDao.GetByEmail(f.Email)
	if err == nil {
		token, err := generateRandomString(32)

		now := time.Now()
		exp := now.Add(4 * time.Hour)
		user.LostPasswordToken = token
		user.LostPasswordTokenExpiration = &exp

		_, err = userDao.UpsertByID("", user)
		if err != nil {
			logger.WithField("error", err).WithField("user", *user).Warn("unable to update user")
			helpers.SendJSONError(w, "Error while update user", http.StatusInternalServerError)
			return
		}

		templateData := struct {
			Name  string
			Email string
			Token string
		}{
			Name:  user.FirstName,
			Email: user.Email,
			Token: user.LostPasswordToken,
		}

		err = helpers.SendEmail("auth/forgot_password", templateData, user.Email, user.Locale)
		if err != nil {
			logger.WithField("error", err).WithField("user", *user).Warn("unable to send forgot password email")
			helpers.SendJSONError(w, "Error while send email", http.StatusInternalServerError)
			return
		}

		logger.WithField("email", user.Email).WithField("token", user.LostPasswordToken).WithField("exp", exp).Debug("user lost password token generated & sended")
	}

	helpers.SendJSONOk(w, nil)
	return
}

func ForgotPasswordConfirmController(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var f CheckTokenForm
	err := helpers.GetJSONContent(&f, r)
	if err != nil {
		logger.WithField("error", err).WithField("form", f).Warn("ForgotPasswordConfirmController error")
		return
	}

	userDao, err := dao.GetUserDAO(dao.DAOMongo)
	if err != nil {
		logger.WithField("error", err).Fatal("unable to connect to mongo db")
		helpers.SendJSONError(w, "Error while sending forgot password token", http.StatusInternalServerError)
		return
	}

	user, err := userDao.GetByEmail(f.Email)
	if err == nil {
		now := time.Now()
		exp := user.LostPasswordTokenExpiration
		if exp != nil && exp.After(now) && user.LostPasswordToken == f.Token {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(f.Password), bcrypt.DefaultCost)
			if err != nil {
				logger.WithField("error", err).Fatal("unable to generate hash password")
				helpers.SendJSONError(w, "Error while generate hash password", http.StatusInternalServerError)
				return
			}
			user.Password = string(hashedPassword[:])

			_, err = userDao.UpsertByID("", user)
			if err != nil {
				logger.WithField("error", err).WithField("user", *user).Warn("unable to update user")
				helpers.SendJSONError(w, "Error while update user", http.StatusInternalServerError)
				return
			}
		} else {
			helpers.SendJSONError(w, "Token for reset password is expired or unknown.", http.StatusUnauthorized)
			return
		}
	}

	helpers.SendJSONOk(w, nil)
	return
}
