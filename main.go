package main

import (
	"fmt"
	"net/http"
	"encoding/json"
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
		var url Url

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&url)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		for i := 0; i < 1000; i++ {

			shortUrl := generateUniqueCode()
			UrlMappings[shortUrl] = url.Url
			fmt.Println(shortUrl)
		}
		

	})


	// starting server on port 8080
	fmt.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}



}




