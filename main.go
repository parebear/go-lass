package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"net/url"
)
// struct to read the url json
type Url struct {
	Url		string		`json:"url"`
}
var UrlMappings = make(map[string]string)


func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World")
	})

	http.HandleFunc("/{code}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprintf(w, "You requested: %s\n", r.URL.Path)
	})

	http.HandleFunc("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		var rawUrl Url

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&rawUrl)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		parsedURL, err := url.Parse(rawUrl.Url)
		if err != nil {
			fmt.Println("Error parsing URL:", err)
			return
		}
		if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
			fmt.Println("Proper scheme not provided by URL")
			fmt.Println("", parsedURL.Scheme)
			return
		}
		if parsedURL.Host == "" || parsedURL.Host == "localhost" || parsedURL.Host == "127.0.0.1" {
			fmt.Println("Host was either empty or not allowed")
			return
		}

		shortUrl := generateUniqueCode()
		UrlMappings[shortUrl] = parsedURL.String()

		fmt.Println("",parsedURL.String(),shortUrl)
		

	})


	// starting server on port 8080
	fmt.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}



}




