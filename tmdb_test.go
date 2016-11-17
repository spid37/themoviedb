package themoviedb_test

import (
	"log"
	"os"
	"testing"

	"github.com/spid37/themoviedb"
)

var apiKey string

var movieID int64 = 293660

func init() {
	apiKey = os.Getenv("TMDB_KEY")
	if apiKey == "" {
		log.Fatal("No apiKey set")
	}
}

func TestPopular(t *testing.T) {
	tmdb := themoviedb.NewTmdb(apiKey)
	_, err := tmdb.Popular()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSearch(t *testing.T) {
	tmdb := themoviedb.NewTmdb(apiKey)
	_, err := tmdb.Search("deadpool")
	if err != nil {
		t.Fatal(err)
	}
}

func TestNoSearch(t *testing.T) {
	tmdb := themoviedb.NewTmdb(apiKey)
	_, err := tmdb.Search("")
	// should return err "query must be provided"
	if err != nil {
		if err.Error() != "query must be provided" {
			t.Fatal(err)
		}
	} else {
		t.Fatal("query is empty! i should not be here!")
	}
}

func TestMovie(t *testing.T) {
	tmdb := themoviedb.NewTmdb(apiKey)
	_, err := tmdb.Movie(movieID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMovieWrongKey(t *testing.T) {
	tmdb := themoviedb.NewTmdb("wrongkey")
	_, err := tmdb.Movie(movieID)
	// should return error "Invalid API key: You must be granted a valid key."
	if err != nil {
		if err.Error() != "Invalid API key: You must be granted a valid key." {
			t.Fatal(err)
		}
	} else {
		t.Fatal("Api key is invalid! i should not be here!")
	}
}
