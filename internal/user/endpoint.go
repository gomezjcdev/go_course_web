package user

import (
	"encoding/json"
	"github.com/gomezjcdev/go_course_web/pkg/meta"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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

	UpdateReq struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"`
		Err    string      `json:"error,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
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
			json.NewEncoder(w).Encode(Response{Status: 400, Err: "Invalid request format"})
			return
		}

		if req.FirstName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: "first name is required"})
			return
		}

		if req.LastName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: "last name is required"})
			return
		}

		user, createErr := s.Create(req.FirstName, req.LastName, req.Email, req.Phone)
		if createErr != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: createErr.Error()})
			return
		}
		json.NewEncoder(w).Encode(Response{Status: 200, Data: user})
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		path := mux.Vars(r)

		user, err := s.Get(path["id"])

		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(Response{Status: 200, Data: user})
	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		v := r.URL.Query()

		filters := Filters{
			FirstName: v.Get("first_name"),
			LastName:  v.Get("last_name"),
		}

		limit, _ := strconv.Atoi(v.Get("limit"))
		page, _ := strconv.Atoi(v.Get("page"))

		count, err := s.Count(filters)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(Response{Status: 500, Err: err.Error()})
			return
		}

		metaRes, err := meta.New(page, limit, count)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(Response{Status: 500, Err: err.Error()})
			return
		}

		users, err := s.GetAll(filters, metaRes.Offset(), metaRes.Limit())
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(Response{Status: 200, Data: users, Meta: metaRes})
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		var req UpdateReq

		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: "Invalid request format"})
			return
		}

		if req.FirstName != nil && *req.FirstName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: "first name is required"})
			return
		}

		if req.LastName != nil && *req.LastName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: "last name is required"})
			return
		}

		path := mux.Vars(r)
		id := path["id"]

		updateErr := s.Update(id, req.FirstName, req.LastName, req.Email, req.Phone)

		if updateErr != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(Response{Status: 404, Err: "user doesn't exist"})
			return
		}

		json.NewEncoder(w).Encode(Response{Status: 200, Data: "success"})
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)

		err := s.Delete(path["id"])

		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: err.Error()})
		}

		json.NewEncoder(w).Encode(Response{Status: 200, Data: "success"})
	}
}
