package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Name string `json:"name"`
}
var userCache = make(map[int]User)

func main(){
	logServe := http.NewServeMux()
	logServe.HandleFunc("/", logServRoot)

	logServe.HandleFunc("POST /users", createUser)

	fmt.Println("Server listening to :8080")
	http.ListenAndServe(":8080", logServe)
}

func logServRoot(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Hello"))
}

func createUser(w http.ResponseWriter, r *http.Request){
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	userCache[len(userCache)+1] = user
	w.WriteHeader(http.StatusNoContent)
}