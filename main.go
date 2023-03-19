package main

import (
	"github.com/gomezjcdev/go_course_web/internal/user"
	"github.com/gomezjcdev/go_course_web/pkg/bootstrap"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"
)

func main() {

	router := mux.NewRouter()
	_ = godotenv.Load()
	logger := bootstrap.InitLogger()

	db, err := bootstrap.DBConnection()

	if err != nil {
		log.Fatal(err)
	}

	userRepo := user.NewRepo(logger, db)
	userSrv := user.NewService(logger, userRepo)
	userEnd := user.MakeEndpoints(userSrv)

	router.HandleFunc("/users/{id}", userEnd.Get).Methods("GET")
	router.HandleFunc("/users", userEnd.GetAll).Methods("GET")
	router.HandleFunc("/users", userEnd.Create).Methods("POST")
	router.HandleFunc("/users/{id}", userEnd.Update).Methods("PATCH")
	router.HandleFunc("/users/{id}", userEnd.Delete).Methods("DELETE")

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
