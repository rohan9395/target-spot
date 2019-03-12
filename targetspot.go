 package main

 import (
	 "fmt"
	 "github.com/target-spot/config"
	 "github.com/target-spot/controller"
	 "log"
	 "net/http"
	 "os"
 )

 func main() {
	 ready := true
	 var endpointMap map[string]util.Endpoint

	 err := util.ReadMockEndpointsData(&endpointMap)

 	fmt.Println("listening...")

	 router := controller.GetRouter(endpointMap, &ready)
	 err = http.ListenAndServe(":"+os.Getenv("PORT"), router)

 	if err != nil {
 		log.Fatal("ListenAndServe:", err)
 	}
 }

 //func GetPort() string {
	// var port = os.Getenv("PORT")
	// // Set a default port if there is nothing in the environment
	// if port == "" {
	//	 port = "4747"
	//	 fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	// }
	// return ":" + port
 //}
