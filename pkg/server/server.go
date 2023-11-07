package server


import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/segunjkf/server/pkg/database"
)

// server is a http server
type Server struct {
	ctx	context.Context
	db  database.Database
}


func New(ctx context.Context, db database.Database) *Server{
	return &Server{
		ctx: ctx,
		db: db,
	}
}

// HandleFuncCreateUser handles the path /user/create
// create -> post
func (s *Server) HandleCreateUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost, http.MethodPut:
		// Check that the input is JSON.
		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return 
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Could not read request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Unmarshal the body.
		var u database.User
		err = json.Unmarshal(body, &u)
		if err != nil {
			log.Printf("Could not unmarshal request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Validation:
		// 1. User Name should not be empty.
		// 2. User must not exist in order to be created.
		if u.Name == "" {
			log.Print("Empty username")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		got := s.db.GetUser(s.ctx, u.Name)
		if got != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("User already exists: %v", u.Name)))
			return
		}

		log.Printf("Create User: %v", u.Name)
		// Write to database.
		err = s.db.Create(s.ctx, u)
		if err != nil {
			log.Printf("Could not create user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}


// func (s *server) HandleFuncHome(w http.ResponseWriter, r *http.Request) {
// 	getHtml, err := ioutil.ReadFile("static/index.html")
// 	if err != nil {
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Add("Content-Type", "text/html")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(getHtml)
// }

 // HandleReadUsers handles the /user request for fetching users
func (s *Server) HandleUsers(w http.ResponseWriter, r *http.Request) {
    // Fetch the name from the params. Common for all methods of this route.
	params := mux.Vars(r)
	name := params["name"]

	switch r.Method {
	case http.MethodGet:
		
		user :=  s.db.GetUser(s.ctx, name)
		if user == nil {
			log.Print("Could not get user")
			w.WriteHeader(http.StatusNotFound) // 500
			return
		}

		msg, err := json.Marshal(user)
		if err != nil {
			log.Printf("Could not read the request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError) // HTTP 500
			return
		}
		log.Printf("Getting user %v", user.Name)

		w.Header().Add("Content-Type", "application/json")
		w.Write(msg)
	case http.MethodPatch:
		// partial update
		// Check if the input is json

		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}  
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Could not read the request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError) // HTTP 500
			return
		}
		defer r.Body.Close()

		// unmarshall body
		var updatedUser database.User
		err = json.Unmarshal(body, &updatedUser)
		if err != nil {
			log.Printf("Error unmarshal request body: %v", err)
			w.WriteHeader(http.StatusBadRequest) // HTTP 400
			return
		}

		log.Printf("updating user %v", name)

		userInfo := s.users[name]
		if updatedUser.Age != 0{
			userInfo.age = updatedUser.Age
		}
		if updatedUser.Email != ""{
			userInfo.email = updatedUser.Email
		}
		s.users[name] = userInfo
// 	case http.MethodDelete:
// 		log.Printf("deleting user: %s", name) 
// 		delete(s.users, name)	
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
 	}

}
