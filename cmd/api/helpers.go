package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/nazevedo3/image_share/internal/data"
)

// Define an clientResponse type.
type clientResponse map[string]interface{}

// writeJSON is a helper that writes the response to the client in JSON format.
func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

// extractImageInfo unmarshals the image data in the payload of the request from the frontend
func (app *application) extractImageInfo(imageInfo string) (*data.Image, error) {
	image := data.Image{}
	err := json.Unmarshal([]byte(imageInfo), &image)
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// deleteFile removes the stored image from the file system
func (app *application) deleteFile(filePath string) error {
	err := os.Remove(filePath[22:])
	if err != nil {
		return err
	}
	return nil
}

// writeJSON is a helper that writes the response to the client in JSON format.
func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	return nil
}
