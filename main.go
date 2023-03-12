package main

import (
	"fmt"
	"github.com/gomezjcdev/go_course_web/internal/user"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	router := mux.NewRouter()
	_ = godotenv.Load()
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	dsn := fmt.Sprintf("host=%s user=%s password='%s' dbname=%s port=%s sslmode=disable TimeZone=America/Bogota",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_PORT"),
	)

	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	db.AutoMigrate(&user.User{})

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

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
