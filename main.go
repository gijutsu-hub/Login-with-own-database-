package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type ErrorTemplateData struct {
	ErrorMessage string
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var errorTemplate = template.Must(template.ParseFiles("templates/login.html"))

func main() {
	router := mux.NewRouter()

	// Define a dynamic route with a URL parameter
	router.HandleFunc("/api/{id}", authorize(userHandler)).Methods("GET")
	router.HandleFunc("/", dbgo).Methods("GET", "POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Middleware function to check for authorization token in request headers
func authorize(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "Sumon" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// Handler function to handle authorized requests
func userHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprintf(w, "User ID: %s\n", id)
}

func binarySearchCredentials(creds []Credentials, username string, password string) bool {
	low := 0
	high := len(creds) - 1

	for low <= high {
		mid := (low + high) / 2
		if creds[mid].Username < username {
			low = mid + 1
		} else if creds[mid].Username > username {
			high = mid - 1
		} else {
			// Username found, check password
			if creds[mid].Password == password {
				return true
			} else {
				return false
			}
		}
	}

	// Username not found
	return false
}

func dbgo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Execute the template
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Open the JSON file and decode the data into a slice of Credentials structs
		file, err := os.Open("./database/web-client/login.json")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		var creds []Credentials
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&creds)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Check if the input credentials match any of the data in the JSON file
		if !binarySearchCredentials(creds, username, password) {
			errorTemplateData := ErrorTemplateData{ErrorMessage: "Invalid credentials"}
			w.WriteHeader(http.StatusUnauthorized)
			err := errorTemplate.Execute(w, errorTemplateData)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			return
		}
		fmt.Print("yes")
		// Render success page if the credentials are valid

	}
}
