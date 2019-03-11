 package main

 import (
 	"fmt"
 	"log"
 	"net/http"
 	"os"
 )

 func main() {
 	http.HandleFunc("/", handler)
 	fmt.Println("listening...")
 	err := http.ListenAndServe(GetPort(), nil)
 	if err != nil {
 		log.Fatal("ListenAndServe: ", err)
 	}
 }

 func handler(w http.ResponseWriter, r *http.Request) {
 	fmt.Fprintf(w, "Rohans App!")
 }

 // Get the Port from the environment so we can run on Heroku
func GetPort() string {
var port = os.Getenv("PORT")
// Set a default port if there is nothing in the environment
if port == "" {
	port = "4747"
	fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
}
	return ":" + port
}
