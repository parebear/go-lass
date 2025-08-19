package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"net/url"
	"strings"
	"sync"
	_ "embed"
)
// struct to read the url json
type ShortenRequest struct {
	Url		string		`json:"url"`
}

type ShortenResponse struct {
	ShortURL		string		`json:"short_url"`
}
var BASE_URL string = "http://localhost:8080"
//go:embed index.html
var indexHTML string
var (
	UrlMappings = make(map[string]string)
	mapMutex = &sync.RWMutex{}
)


func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Handler hit: Path=%s Method=%s\n", r.URL.Path, r.Method)
		if r.URL.Path == "/" &&r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(indexHTML))
			return
		}
	})


	http.HandleFunc("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		var rawUrl ShortenRequest

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
		mapMutex.Lock()
		defer mapMutex.Unlock()
		UrlMappings[shortUrl] = parsedURL.String()
		

		shortenedUrl := ShortenResponse{ShortURL: BASE_URL + "/" + shortUrl}
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(shortenedUrl); err != nil {
			fmt.Println("Error encoding JSON:", err)
		}


	})

	http.HandleFunc("/{code}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}

		if r.URL.Path == "/" || strings.HasPrefix(r.URL.Path, "/api/") {
			 http.NotFound(w, r)
			 return
		}

		shortCode := r.URL.Path[1:]
		mapMutex.RLock()
		defer mapMutex.RUnlock()


		mapMutex.RLock()
		defer mapMutex.RUnlock()
		if originalURL, exists := UrlMappings[shortCode]; !exists {
			http.NotFound(w, r)
		} else {
			http.Redirect(w, r, originalURL, http.StatusFound)
			return
		}





	})


	// starting server on port 8080
	fmt.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}



}




