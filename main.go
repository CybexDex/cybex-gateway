package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func main() {
	router := mux.NewRouter()

	// router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	// router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	// router.HandleFunc("/api/contacts/new", controllers.CreateContact).Methods("POST")
	// router.HandleFunc("/api/me/contacts", controllers.GetContactsFor).Methods("GET") //  user/2/contacts

	// router.Use(app.JwtAuthentication) //attach JWT auth middleware
	router.Use(loggingMiddleware)

	//router.NotFoundHandler = app.NotFoundHandler

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}

	// TODO: graceful shutdown
}
