package main

import "fmt"
import "net/http"

func main() {
	fmt.Println("hello world!")

	resp, err := http.Get("http://example.com")
	if err != nil {
		fmt.Println("Error: This didn't work")
	}
	fmt.Println(resp)
}
