package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jpschew/GoSchool-Projects/project-5-build-restful-api/courseapi/api"
	"log"
	"net/http"
)

func main() {

	// declare a new mux router
	router := mux.NewRouter()

	// route and handler functions for course and module
	// without specifying the methods using .Method()
	// it will accept all methods
	router.HandleFunc("/api/v1/", api.Home).Methods("GET")
	router.HandleFunc("/api/v1/courses", api.Allcourses).Methods("GET")
	router.HandleFunc("/api/v1/courses/{courseid}", api.Course).Methods("GET", "DELETE", "POST", "PUT")
	router.HandleFunc("/api/v1/courses/{courseid}/{moduleid}", api.Course).Methods("GET", "DELETE", "POST", "PUT")

	// route and handler functions for api key
	router.HandleFunc("/api/v1/genkey", api.GenKey).Methods("GET")
	router.HandleFunc("/api/v1/addupdatekey", api.AddUpdateKey).Methods("POST", "PUT")
	router.HandleFunc("/api/v1/deletekey", api.DeleteKey).Methods("POST")

	// listen at port 5000
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
