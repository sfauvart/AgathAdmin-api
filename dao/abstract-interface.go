package dao

// UserDAO is the DAO interface to work with users
type BaseDAO interface {

	// GetUserByID returns a user by its ID
	GetByID(ID string) (interface{}, error)

	// GetAllUsers returns all users with paging capability
	GetAll(start, end int, sort string) (interface{}, int, error)

	// SaveUser saves the user
	Save(model interface{}) error

	// UpsertUser updates or creates a user
	UpsertByID(ID string, model interface{}) (bool, error)

	// DeleteUser deletes a users by its ID
	DeleteByID(ID string) error
}

type AbstractDAO interface {
	BaseDAO
}
