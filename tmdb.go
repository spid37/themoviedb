package themoviedb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const baseURL = "https://api.themoviedb.org/3"
const posterBaseURL = "https://image.tmdb.org/t/p/w500/"
const backdropBaseURL = "https://image.tmdb.org/t/p/w780/"

// NewTmdb get a new tmdb instance
func NewTmdb(apiKey string) *TMDB {
	tmdb := new(TMDB)
	tmdb.APIKey = apiKey
	return tmdb
}

// TMDB -
type TMDB struct {
	APIKey string
}

// SearchResult result of a search query
type SearchResult struct {
	Page         int     `json:"page"`
	Results      []Movie `json:"results"`
	TotalResults int     `json:"total_results"`
	TotalPages   int     `json:"total_pages"`
}

// Movie a movie entity
type Movie struct {
	ID            int64   `json:"id"`
	ImdbID        string  `json:"imdb_id"`
	Title         string  `json:"title"`
	OriginalTitle string  `json:"original_title"`
	Tagline       string  `json:"tagline"`
	Overview      string  `json:"overview"`
	Genres        []Genre `json:"genres"`
	Runtime       int     `json:"runtime"`
	ReleaseDate   string  `json:"release_date"`
	Backdrop      string  `json:"backdrop_path"`
	Poster        string  `json:"poster_path"`
	Adult         bool    `json:"adult"`
}

// Genre of the movie
type Genre struct {
	Name string `json:"name"`
}

// GetPosterURL get the url for the poster image
func (m *Movie) GetPosterURL() string {
	if m.Poster == "" {
		return ""
	}
	return fmt.Sprintf("%s%s", posterBaseURL, m.Poster)
}

// GetBackdropURL get the url for the poster image
func (m *Movie) GetBackdropURL() string {
	if m.Backdrop == "" {
		return ""
	}
	return fmt.Sprintf("%s%s", backdropBaseURL, m.Backdrop)
}

// Movie gte a movie by its ID
func (t *TMDB) Movie(ID int64) (*Movie, error) {
	var err error
	movie := new(Movie)

	// build the request URL
	u, err := url.Parse(baseURL)
	if err != nil {
		return movie, err
	}
	u.Path = fmt.Sprintf("%s%s/%d", u.Path, "/movie", ID)
	q := u.Query()
	q.Set("api_key", t.APIKey)
	u.RawQuery = q.Encode()

	// create http client with timout
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	// make the request
	response, err := netClient.Get(u.String())
	if err != nil {
		return movie, err
	}
	err = processResponse(response, movie)

	return movie, err
}

// Popular get popular movies
func (t *TMDB) Popular() (*SearchResult, error) {
	var err error
	result := new(SearchResult)

	// build the request URL
	u, err := url.Parse(baseURL)
	if err != nil {
		return result, err
	}
	u.Path = fmt.Sprintf("%s%s", u.Path, "/movie/popular")
	q := u.Query()
	q.Set("api_key", t.APIKey)
	u.RawQuery = q.Encode()

	// create http client with timout
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	// make the request
	response, err := netClient.Get(u.String())
	if err != nil {
		return result, err
	}
	err = processResponse(response, result)

	return result, err
}

// Search for a movie
func (t *TMDB) Search(query string) (*SearchResult, error) {
	var err error
	results := new(SearchResult)

	// build the request URL
	u, err := url.Parse(baseURL)
	if err != nil {
		return results, err
	}
	u.Path = fmt.Sprintf("%s%s", u.Path, "/search/movie")
	q := u.Query()
	q.Set("api_key", t.APIKey)
	q.Set("query", query)
	q.Set("include_adult", "false")
	u.RawQuery = q.Encode()

	// create http client with timout
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	// make the request
	response, err := netClient.Get(u.String())
	err = processResponse(response, results)

	return results, err
}

func processResponse(response *http.Response, result interface{}) error {
	var err error
	// close the body later
	defer response.Body.Close()

	// all success responses should be 200
	if response.StatusCode != 200 {
		err = getError(response)
	}
	if err != nil {
		return err
	}
	// no errors found
	// close the request body
	// add response to SearchResult
	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(buf, result)
}

func getError(response *http.Response) error {
	var err error
	responseErr := new(ResponseError)

	contentType := response.Header.Get("Content-Type")
	if contentType != "application/json;charset=utf-8" && contentType != "application/json" {
		err = fmt.Errorf("Invalid response status code (%d)", response.StatusCode)
		return err
	}

	// add response to SearchResult
	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, &responseErr)
	if err != nil {
		return err
	}

	if responseErr.Message == "" && len(responseErr.Errors) > 0 {
		responseErr.Message = strings.Join(responseErr.Errors, ",")
	}

	return errors.New(responseErr.Message)
}

// ResponseError error response from tmdb api
type ResponseError struct {
	Code    int      `json:"status_code"`
	Message string   `json:"status_message"`
	Errors  []string `json:"errors"`
}
