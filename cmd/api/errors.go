package main

import (
	"fmt"
	"net/http"
)

// logError is a helper method for logging an error message
func (app *application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

// errorResponse is used for sending JSON formatted error messages to the client.
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	resp := clientResponse{"error": message}
	// Write the response using the  helper method.
	err := app.writeJSON(w, status, resp)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// serverErrorResponse will be used to send a 500 Internal Server Error status code and
// JSON response to the client.
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// notFoundResponse will be used to send a 404 Not Found status code and
// JSON response to the client.
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// methodNotAllowedResponse will be used to send a 405 Method Not Allowed
// status code and JSON response to the client.
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}
