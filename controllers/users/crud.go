package users

import (
	"net/http"
	"github.com/sfauvart/Agathadmin-api/dao"
  "github.com/sfauvart/Agathadmin-api/helpers"
	"github.com/sfauvart/Agathadmin-api/models"
	"gopkg.in/mgo.v2"
  logger "github.com/Sirupsen/logrus"
  "strconv"
)

func GetAll(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  startStr := helpers.ParamAsString("start", r)
	endStr := helpers.ParamAsString("end", r)

	start := dao.NoPaging
	end := dao.NoPaging
	var err error
	if startStr != "" && endStr != "" {
		start, err = strconv.Atoi(startStr)
		if err != nil {
			start = dao.NoPaging
		}
		end, err = strconv.Atoi(endStr)
		if err != nil {
			end = dao.NoPaging
		}
	}

  userDao, err := dao.GetUserDAO(dao.DAOMongo)
	if err != nil {
		logger.WithField("error", err).Fatal("unable to connect to mongo db")
		w.WriteHeader(500)
		return
	}

	// Fetch users
	users, err := userDao.GetAll(start, end)
  if err != nil {
		logger.WithField("error", err).Warn("unable to retrieve users")
		helpers.SendJSONError(w, "Error while retrieving users", http.StatusInternalServerError)
		return
	}

	logger.WithField("users", users).Debug("users found")
	helpers.SendJSONOk(w, users)
}

// Get retrieve an entity by id
func Get(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// get the user ID from the URL
	userID := helpers.ParamAsString("id", r)

	userDao, err := dao.GetUserDAO(dao.DAOMongo)
	if err != nil {
		logger.WithField("error", err).Fatal("unable to connect to mongo db")
		w.WriteHeader(500)
		return
	}

	// find user
	user, err := userDao.GetByID(userID)
	if err != nil {
		if err == mgo.ErrNotFound {
			logger.WithField("error", err).WithField("user ID", userID).Warn("unable to retrieve user by ID")
			helpers.SendJSONNotFound(w)
			return
		}

		logger.WithField("error", err).WithField("user ID", userID).Warn("unable to retrieve user by ID")
		helpers.SendJSONError(w, "Error while retrieving user by ID", http.StatusInternalServerError)
		return
	}

	logger.WithField("users", user).Debug("user found")
	helpers.SendJSONOk(w, user)
}

// Create create an entity
func Create(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// user to be created
	user := &models.User{}
	// get the content body
	err := helpers.GetJSONContent(user, r)

	if err != nil {
		logger.WithField("error", err).Warn("unable to decode user to create")
		helpers.SendJSONError(w, "Error while decoding user to create", http.StatusBadRequest)
		return
	}

	userDao, err := dao.GetUserDAO(dao.DAOMongo)
	if err != nil {
		logger.WithField("error", err).Fatal("unable to connect to mongo db")
		helpers.SendJSONError(w, "Error while creating user", http.StatusInternalServerError)
		return
	}

	// save user
	err = userDao.Save(user)
	if err != nil {
		logger.WithField("error", err).WithField("user", *user).Warn("unable to create user")
		helpers.SendJSONError(w, "Error while creating user", http.StatusInternalServerError)
		return
	}

	// send response
	helpers.SendJSONOk(w, user)
}

// Update update an entity by id
func Update(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// get the user ID from the URL
	userID := helpers.ParamAsString("id", r)

	// user to be created
	user := &models.User{}
	// get the content body
	err := helpers.GetJSONContent(user, r)

	if err != nil {
		logger.WithField("error", err).Warn("unable to decode user to create")
		helpers.SendJSONError(w, "Error while decoding user to create", http.StatusBadRequest)
		return
	}

	userDao, err := dao.GetUserDAO(dao.DAOMongo)
	if err != nil {
		logger.WithField("error", err).Fatal("unable to connect to mongo db")
		helpers.SendJSONError(w, "Error while creating user", http.StatusInternalServerError)
		return
	}

	// save user
	_, err = userDao.UpsertByID(userID, user)
	if err != nil {
		logger.WithField("error", err).WithField("user", *user).Warn("unable to create user")
		helpers.SendJSONError(w, "Error while creating user", http.StatusInternalServerError)
		return
	}

	// send response
	helpers.SendJSONOk(w, user)
}

// Delete delete an entity by id
func Delete(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// get the user ID from the URL
	userID := helpers.ParamAsString("id", r)

	userDao, err := dao.GetUserDAO(dao.DAOMongo)
	if err != nil {
		logger.WithField("error", err).Fatal("unable to connect to mongo db")
		helpers.SendJSONError(w, "Error while deleting user by ID", http.StatusInternalServerError)
		return
	}

	// find user
	err = userDao.DeleteByID(userID)
	if err != nil {
		logger.WithField("error", err).WithField("user ID", userID).Warn("unable to delete user by ID")
		helpers.SendJSONError(w, "Error while deleting user by ID", http.StatusInternalServerError)
		return
	}

	logger.WithField("userID", userID).Debug("user deleted")
	helpers.SendJSONOk(w, nil)
}
