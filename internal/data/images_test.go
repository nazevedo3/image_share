package data

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var testdata = []Image{
	{
		Title:       "Gopher",
		Description: "Go's Mascot",
		Tags:        []string{"golang", "gopher", "codingIsAwesome"},
		Path:        "https://locahost:8080/upload-432049540.jpeg",
	}, {
		Title:       "Coffee",
		Description: "Cup of Coffee",
		Tags:        []string{"#coffeeislife", "#coffeefirst"},
		Path:        "https://locahost:8080/upload-9324234.jpeg",
	},
}

func TestImage(t *testing.T) {
	t.Parallel()
	_ = Image{
		ID:          123,
		Title:       "Title for test image",
		Description: "Description for test image",
		Tags:        []string{"tag1", "tag2", "tag3"},
		Path:        "https://locahost:8080/upload-432049540.jpeg",
	}
}

func TestDBActions(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Initialize container and the database
	container, db, err := CreateTestContainer(ctx, "testdb")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer container.Terminate(ctx)

	// use the migration files to seed the database
	mig, err := NewMigrator(db)
	if err != nil {
		t.Fatal(err)
	}

	err = mig.Up()
	if err != nil {
		t.Fatal(err)
	}

	// create ImageModel using the test database
	m := ImageModel{db}

	//Insert mock image into the database
	err = m.InsertImage(&testdata[0])
	if err != nil {
		t.Errorf("want no error inserting images, but got %v", err)
	}
	wantAllImages := 1
	//Retrieve all messages An
	got, err := m.GetAllImages()
	if err != nil {
		t.Errorf("want no error getting all images, but got %v", err)
	}

	if !cmp.Equal(wantAllImages, len(got)) {
		t.Error(!cmp.Equal(wantAllImages, got))
	}
	//Test the Get method returns the expected image
	//Since it's the first entry in the test DB, it should set the ID to 1
	gotImage, err := m.Get(testdata[0].ID)
	if err != nil {
		t.Errorf("want no error getting image by ID, but got %v", err)
	}
	if !cmp.Equal(&testdata[0], gotImage) {
		t.Error(!cmp.Equal(testdata[0], gotImage.ID))
	}
	//Test that Delete doesn't return any errors
	gotNoErr := m.Delete(testdata[0].ID)
	if gotNoErr != nil {
		t.Errorf("want no error deleting image, but got %v", err)
	}
}
