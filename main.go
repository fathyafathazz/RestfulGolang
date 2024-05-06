package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver package
	"github.com/gorilla/mux"
)

// Album struct
type Album struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Year   int    `json:"year"`
}

var db *sql.DB

func main() {
	var err error
	// Open database connection
	db, err = sql.Open("mysql", "root@tcp(localhost:3306)/dbgolang")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize router
	router := mux.NewRouter()


	// Routes
	router.HandleFunc("/albums", getAlbums).Methods("GET")
	router.HandleFunc("/albums/{id}", getAlbum).Methods("GET")
	router.HandleFunc("/albumsPost", createAlbum).Methods("POST")
	router.HandleFunc("/albumsPut/{id}", updateAlbum).Methods("PUT")
	router.HandleFunc("/albumsDelete/{id}", deleteAlbum).Methods("DELETE")
	router.HandleFunc("/albumsCreate", createAlbumViaURL).Methods("POST")
	router.HandleFunc("/albumsUpdate/{id}", updateAlbumViaURL).Methods("PUT")

	// Start server
	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Get all albums
// @Summary Get all albums
// @Description Get all albums from the database
// @Tags albums
// @Produce json
// @Success 200 {array} Album
// @Router /albums [get]
func getAlbums(w http.ResponseWriter, r *http.Request) {
	// Query to get all albums
	rows, err := db.Query("SELECT id, title, artist, year FROM album")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Slice to hold albums
	var albums []Album

	// Iterate over the rows
	for rows.Next() {
		var album Album
		err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Year)
		if err != nil {
			log.Fatal(err)
		}
		albums = append(albums, album)
	}

	// Encode albums to JSON and send response
	json.NewEncoder(w).Encode(albums)
}

// Get single album
// Get single album
// @Summary Get single album by ID
// @Description Get single album from the database by its ID
// @Tags albums
// @Produce json
// @Param id path int true "Album ID"
// @Success 200 {object} Album
// @Router /albums/{id} [get]
func getAlbum(w http.ResponseWriter, r *http.Request) {
	// Get album ID from URL parameters
	params := mux.Vars(r)
	albumID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid album ID", http.StatusBadRequest)
		return
	}

	// Query to get single album
	row := db.QueryRow("SELECT id, title, artist, year FROM album WHERE id = ?", albumID)

	// Initialize album struct to hold data
	var album Album

	// Scan row data into album struct
	err = row.Scan(&album.ID, &album.Title, &album.Artist, &album.Year)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Album not found", http.StatusNotFound)
			return
		}
		log.Fatal(err)
	}

	// Encode album to JSON and send response
	json.NewEncoder(w).Encode(album)
}

// Create a new album
// @Summary Get all albums
// @Description Get all albums from the database
// @Tags albums
// @Produce json
// @Success 200 {array} Album
// @Router /albumsPost [post]
func createAlbum(w http.ResponseWriter, r *http.Request) {
	// Decode request body into album struct
	var album Album
	err := json.NewDecoder(r.Body).Decode(&album)
	if err != nil {
		log.Fatal(err)
	}

	// Insert new album into database
	result, err := db.Exec("INSERT INTO album (title, artist, year) VALUES (?, ?, ?)", album.Title, album.Artist, album.Year)
	if err != nil {
		log.Fatal(err)
	}

	// Get ID of newly inserted album
	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	// Set ID of album struct
	album.ID = int(id)

	// Encode album to JSON and send response
	json.NewEncoder(w).Encode(album)
}

// Update an album
// @Summary Get all albums
// @Description Get all albums from the database
// @Tags albums
// @Produce json
// @Success 200 {array} Album
// @Router /albumsPut/{id} [put]
func updateAlbum(w http.ResponseWriter, r *http.Request) {
	// Get album ID from URL parameters
	params := mux.Vars(r)
	albumID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	// Decode request body into album struct
	var album Album
	err = json.NewDecoder(r.Body).Decode(&album)
	if err != nil {
		log.Fatal(err)
	}

	// Update album in database
	_, err = db.Exec("UPDATE album SET title = ?, artist = ?, year = ? WHERE id = ?", album.Title, album.Artist, album.Year, albumID)
	if err != nil {
		log.Fatal(err)
	}

	// Encode album to JSON and send response
	json.NewEncoder(w).Encode(album)
}

// Delete an album
// @Summary Get all albums
// @Description Get all albums from the database
// @Tags albums
// @Produce json
// @Success 200 {array} Album
// @Router /albumsDelete/{id} [delete]
func deleteAlbum(w http.ResponseWriter, r *http.Request) {
	// Get album ID from URL parameters
	params := mux.Vars(r)
	albumID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	// Delete album from database
	_, err = db.Exec("DELETE FROM album WHERE id = ?", albumID)
	if err != nil {
		log.Fatal(err)
	}

	// Send success message
	json.NewEncoder(w).Encode(map[string]string{"message": "Album deleted successfully"})
}

// Create a new album via URL
// @Summary Get all albums
// @Description Get all albums from the database
// @Tags albums
// @Produce json
// @Success 200 {array} Album
// @Router /albumsCreate [post]
func createAlbumViaURL(w http.ResponseWriter, r *http.Request) {
	// Get parameters from URL query
	title := r.URL.Query().Get("title")
	artist := r.URL.Query().Get("artist")
	yearStr := r.URL.Query().Get("year")

	// Convert year to int
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	// Insert new album into database
	result, err := db.Exec("INSERT INTO album (title, artist, year) VALUES (?, ?, ?)", title, artist, year)
	if err != nil {
		log.Fatal(err)
	}

	// Get ID of newly inserted album
	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare response JSON
	album := Album{
		ID:     int(id),
		Title:  title,
		Artist: artist,
		Year:   year,
	}

	// Encode album to JSON and send response
	json.NewEncoder(w).Encode(album)
}

// Update an album via URL
// @Summary Get all albums
// @Description Get all albums from the database
// @Tags albums
// @Produce json
// @Success 200 {array} Album
// @Router /albumsUpdate/{id} [put]
func updateAlbumViaURL(w http.ResponseWriter, r *http.Request) {
	// Get parameters from URL query
	title := r.URL.Query().Get("title")
	artist := r.URL.Query().Get("artist")
	yearStr := r.URL.Query().Get("year")

	// Convert year to int
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	// Get album ID from URL parameters
	params := mux.Vars(r)
	albumID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	// Update album in database
	_, err = db.Exec("UPDATE album SET title = ?, artist = ?, year = ? WHERE id = ?", title, artist, year, albumID)
	if err != nil {
		log.Fatal(err)
	}

	// Prepare response JSON
	album := Album{
		ID:     albumID,
		Title:  title,
		Artist: artist,
		Year:   year,
	}

	// Encode album to JSON and send response
	json.NewEncoder(w).Encode(album)
}
