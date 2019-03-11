 package main

 import (
 	"fmt"
 	"log"
 	"net/http"
 	"os"
  "github.com/gin-gonic/gin"

 )

 func main() {
 	fmt.Println("listening...")
  router := controller.GetRouter(endpointMap, &ready)
 	err := http.ListenAndServe(GetPort(), router)
 	if err != nil {
 		log.Fatal("ListenAndServe: ", err)
 	}
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
