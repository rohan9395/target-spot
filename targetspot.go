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
 		log.Fatal("ListenAndServe: ", err)
 	}
 }
