package Router

import (
	"log"
	"net/http"

	//this is our user defined handler function to handle the routes
	"test3/packages/Handler"

	"github.com/gorilla/mux"
)

//this function will return the router of the gorrila mux
func StartServer() *mux.Router {

	//declaring the gorilla mux to make use of our routes
	r := mux.NewRouter()
	log.Println("Starting the application")

	//if it is the route to get all the documnets
	r.HandleFunc("/todo", Handler.RootHandleFunc).Methods("GET")

	//route to get the particular document by id
	r.HandleFunc("/todo/{id}", Handler.GetHandlerFunction).Methods("GET")

	//route to add new todo item in the list
	r.HandleFunc("/todo", Handler.AddHandlerFunc).Methods("POST")

	//route to delete todo item from the list
	r.HandleFunc("/todo/{id}", Handler.DeleteHandlerFunc).Methods("DELETE")

	//route to update the data of the todo list
	r.HandleFunc("/todo/{id}", Handler.UpdateHandlerFunc).Methods("PUT")

	http.Handle("/", r)

	//log the error if the server is already in use
	log.Fatal(http.ListenAndServe(":8080", nil))
	return r //returning the route
}
