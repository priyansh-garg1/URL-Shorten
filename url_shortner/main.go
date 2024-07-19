package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

var urlDB = make(map[string]URL)

func generateShortURL(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL))
	dataHasher := hasher.Sum(nil)
	hash := hex.EncodeToString(dataHasher)
	fmt.Println(hash)
	return hash[:8]
}

func createURL(OriginalURL string) string {
	shortURL := generateShortURL(OriginalURL)
	id := shortURL
	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  OriginalURL,
		ShortURL:     shortURL,
		CreationDate: time.Now(),
	}
	return shortURL
}

func getURL(shortURL string) (URL, error) {
	url, ok := urlDB[shortURL]
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}

func RootPagrURL(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HOME PAGE")
}

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	shortURL := createURL(data.URL)
	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	fmt.Fprintf(w, response.ShortURL)
}

func RedirectURL(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect")]
	url, err := getURL(string(id))
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)

}

func main() {
	fmt.Println("URL SHORTNER")
	OriginalURL := "https://example.com/hello-world-1"
	generateShortURL(OriginalURL)

	http.HandleFunc("/", RootPagrURL)
	http.HandleFunc("/shorten", ShortURLHandler)

	//Start the http server on PORT 3000
	fmt.Println("Starting http server on PORT 3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
