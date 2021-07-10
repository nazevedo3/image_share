package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/lib/pq"
)

//Image captures all the necessary data for each image
type Image struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Path        string   `json:"path"`
}

//ImageModel is a wrapper for the DB client
type ImageModel struct {
	DB *sql.DB
}

//SetPath builds the string for the path(permalink/URL)
func (i *Image) SetPath(file *os.File) {
	i.Path = fmt.Sprintf("http://localhost:4000/%s", file.Name())
}

//Insert interacts with the database to store a image
func (m ImageModel) InsertImage(image *Image) error {
	query := `
	INSERT INTO images (title, description, tags, path)
	VALUES($1, $2, $3, $4)
	RETURNING id`

	args := []interface{}{image.Title, image.Description, pq.Array(image.Tags), image.Path}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&image.ID)
}

//GetAllImages returns all of the images in the database
func (m ImageModel) GetAllImages() ([]*Image, error) {
	query := `
	SELECT id, title, description, tags, path
	FROM images
	ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	images := []*Image{}

	for rows.Next() {
		var image Image

		err := rows.Scan(
			&image.ID,
			&image.Title,
			&image.Description,
			pq.Array(&image.Tags),
			&image.Path,
		)
		if err != nil {
			return nil, err
		}

		images = append(images, &image)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return images, nil
}

//Get retrieves a image from the database give a id
func (m ImageModel) Get(id int64) (*Image, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
	SELECT id, title, description, tags, path
	FROM images
	WHERE id = $1`

	var image Image

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&image.ID,
		&image.Title,
		&image.Description,
		pq.Array(&image.Tags),
		&image.Path,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &image, nil
}

//Delete removes an image from the database given it's ID
func (m ImageModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
		DELETE FROM images
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil

}
