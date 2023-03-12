package user

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateReq struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	ErrorRes struct {
		Error string `json:"error"`
	}
)

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		var req CreateReq

		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorRes{"Invalid request format"})
			return
		}

		if req.FirstName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorRes{"first name is required"})
			return
		}

		if req.LastName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorRes{"last name is required"})
			return
		}

		user, createErr := s.Create(req.FirstName, req.LastName, req.Email, req.Phone)
		if createErr != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorRes{createErr.Error()})
			return
		}
		json.NewEncoder(w).Encode(user)
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		path := mux.Vars(r)

		user, err := s.Get(path["id"])

		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(ErrorRes{
				Error: err.Error(),
			})
		}

		json.NewEncoder(w).Encode(user)
	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		users, err := s.GetAll()
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorRes{
				Error: err.Error(),
			})
		}

		json.NewEncoder(w).Encode(users)
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Update User")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Delete User")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}
}
