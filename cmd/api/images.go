package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/nazevedo3/image_share/internal/data"
)

// uploadImageHandler stores the received message in the database
func (app *application) uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	// The request from the client comes with multiform headers
	// Parse the Form and set the size limit between 1 and 2MB
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("myFile")

	if err != nil {
		app.logger.Println(err)
		return
	}
	defer file.Close()

	app.logger.Printf("Uploaded File: %+v\n, File Size: %+v\n, MIME Header: %v\n", handler.Filename, handler.Size, handler.Header)

	tempFile, err := ioutil.TempFile("temp-images", fmt.Sprintf("upload-*%s", filepath.Ext(handler.Filename)))
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fmt.Println(tempFile.Name())

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to a temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	app.logger.Println("Successfully Uploaded File")

	//use the helper method to unmarshal the image object received from the front end
	image, err := app.extractImageInfo(r.PostFormValue("data"))
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	//add the path to the image struct
	image.SetPath(tempFile)

	//Insert into the database
	err = app.models.Images.InsertImage(image)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	//writeResponse to the client
	err = app.writeJSON(w, http.StatusCreated, image)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// getAllImagesHander returns all the images currently in the database
func (app *application) getAllImagesHandler(w http.ResponseWriter, r *http.Request) {
	images, err := app.models.Images.GetAllImages()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	//writeResponse to the client
	err = app.writeJSON(w, http.StatusOK, images)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// getImageByIdHandler returns a specific image given it's ID
func (app *application) getImageByIdHandler(w http.ResponseWriter, r *http.Request) {
	//extract the ID parameter for the request
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	//Retrieve the image from the database
	image, err := app.models.Images.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	//writeResponse to the client
	err = app.writeJSON(w, http.StatusOK, image)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteImageHandler will delete an image from the database given it's ID
func (app *application) deleteImageHandler(w http.ResponseWriter, r *http.Request) {
	//extract the ID parameter for the request
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}
	//Get the image info from the request
	//this will be used to delete the image from the file system
	removeImage, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	defer r.Body.Close()
	//Delete from the database
	err = app.models.Images.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

	}
	//Unmarshal the request body
	image, err := app.extractImageInfo(string(removeImage))
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	//Remove the image from the file system
	err = app.deleteFile(image.Path)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	//writeResponse to the client
	err = app.writeJSON(w, http.StatusOK, clientResponse{"message": "image successfully delete"})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
