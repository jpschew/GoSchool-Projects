package main

import (
	"fmt"
	"github.com/jpschew/GoSchool-Projects/project-4-client-server-for-dental-appointment-system-with-security/datatype"
	"github.com/jpschew/GoSchool-Projects/project-4-client-server-for-dental-appointment-system-with-security/server"
	"github.com/jpschew/GoSchool-Projects/project-4-client-server-for-dental-appointment-system-with-security/utils"
	"net/http"
	"runtime"

	"github.com/gorilla/mux"
)

func init() {

	datatype.InitDentist()
	datatype.InitApptHashTable()
	datatype.InitAppt()
}

func main() {

	// 8 phyiscal processor
	runtime.GOMAXPROCS(runtime.NumCPU())
	// fmt.Println(runtime.NumCPU())

	// for panic recovery
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in  main:", err)
		}
	}()

	// server.StartServer()

	// gorilla mux is used here to limit the request method for different routes
	r := mux.NewRouter()
	r.HandleFunc("/", server.Index).Methods(http.MethodGet)
	r.HandleFunc("/signup", server.Signup)
	r.HandleFunc("/login", server.Login)
	r.HandleFunc("/browse", server.Browse)
	r.HandleFunc("/search", server.SearchAppt)
	r.HandleFunc("/make", server.MakeAppt)
	r.HandleFunc("/list", server.List)
	r.HandleFunc("/edit", server.EditAppt)
	r.HandleFunc("/logout", server.Logout).Methods(http.MethodGet)
	r.Handle("/favicon.ico", http.NotFoundHandler())

	// without https
	// utils.ErrorLogging("Client-Server does not establish a connection.", http.ListenAndServe(":8080", r))

	// https is used for secure communications between client and server
	utils.ErrorLogging("Client-Server does not establish a connection.",
		http.ListenAndServeTLS(":8081",
			"key/cert.pem",
			"key/key.pem",
			r))

}
