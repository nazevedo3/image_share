package data

import (
	"database/sql"
	"errors"
)

//ErrRecordNotFound global variable set for when a record in the database can't be found
var (
	ErrRecordNotFound = errors.New("record not found")
)

//Models is a wrapper for the ImageModel
type Models struct {
	Images ImageModel
}

//NewModel initilizes a new Model struct for the db client
func NewModels(db *sql.DB) Models {
	return Models{
		Images: ImageModel{DB: db},
	}
}
