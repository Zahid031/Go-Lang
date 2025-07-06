package main
import (
	"fmt"
	"net/http"
	)

func helloHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Hello,World")
}

func aboutHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "I am Zahid")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/about",aboutHandler)
	fmt.Println("Server is running on port 8000")
	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		fmt.Println("Error starting the server", err)
	}
	
}