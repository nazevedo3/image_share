package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// routes is used to register the handlers
func (app *application) routes() http.Handler {
	router := httprouter.New()
	//settings to allow for Preflight CORS
	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
			header.Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
			header.Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		}
		w.WriteHeader(http.StatusNoContent)
	})
	// Convert the notFoundResponse helper to a http.Handler and
	// the methodNotAllowedResponse to customer error handler for 405 responses.
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/api/images", app.getAllImagesHandler)
	router.HandlerFunc(http.MethodGet, "/api/images/:id", app.getImageByIdHandler)
	router.HandlerFunc(http.MethodPost, "/api/images", app.uploadImageHandler)
	router.HandlerFunc(http.MethodDelete, "/api/images/:id", app.deleteImageHandler)
	//Allows the serving for static assets (images)
	router.ServeFiles("/temp-images/*filepath", http.Dir("temp-images"))

	return app.enableCORS(router)
}
