package users

import (
	logger "github.com/Sirupsen/logrus"
	"github.com/sfauvart/Agathadmin-api/dao"
	"github.com/sfauvart/Agathadmin-api/helpers"
	"github.com/sfauvart/Agathadmin-api/models"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"math"
	"net/http"
	"strconv"
)

func GetAll(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("size")
	sortStr := r.URL.Query().Get("sort")

	page := 1
	size := dao.DefaultElementsPerPage
	var err error
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
	}
	if sizeStr != "" {
		size, err = strconv.Atoi(sizeStr)
	}
	if size > dao.MaxElementsPerPage {
		size = dao.MaxElementsPerPage
	}

	start := (page - 1) * size
	end := (page * size)

	if sortStr == "" {
		sortStr = "+id"
	}

	userDao, err := dao.GetUserDAO(dao.DAOMongo)
	if err != nil {
		logger.WithField("error", err).Fatal("unable to connect to mongo db")
		helpers.SendJSONError(w, "Error while retrieving users", http.StatusInternalServerError)
		return
	}

	// Fetch users
	users, te, err := userDao.GetAll(start, end, sortStr)
	if err != nil {
		logger.WithField("error", err).Warn("unable to retrieve users")
		helpers.SendJSONError(w, "Error while retrieving users", http.StatusInternalServerError)
		return
	}

	pagination := helpers.Pagination{
		CurrentPage:      page,
		TotalPages:       int(math.Ceil(float64(te) / float64(size))),
		NumberOfElements: len(users.([]models.User)),
		TotalElements:    te,
		FirstPage:        true,
		LastPage:         true,
	}
	sort := helpers.Sort{
		Ascending: true,
		Property:  sortStr[1:],
	}
	sortStrDir := sortStr[:1]
	if sortStrDir == "-" {
		sort.Ascending = false
	}
	if page > 1 {
		pagination.FirstPage = false
	}
	if page*size < te {
		pagination.LastPage = false
	}

	json := helpers.PaginatedResult{
		Pagination: pagination,
		Sort:       sort,
		Data:       users,
	}

	logger.WithField("users", users).Debug("users found")
	helpers.SendJSONOk(w, json)
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.WithField("error", err).Fatal("unable to generate hash password")
		helpers.SendJSONError(w, "Error while generate hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword[:])

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
