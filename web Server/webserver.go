package main

//importing net/http to handle the http request and response
//importing io/ioutil to handle the input and output
import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//function to read the file it take a filename as an argument and
//return the string and error as type
func loadFile(fileName string) (string, error) {

	bytes, err := ioutil.ReadFile(fileName)
	//checking if the errorr happens or not
	if err != nil {
		return "", err
	}
	return string(bytes), nil

}

//function to handle the request
func userHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi this is user page")
	//r.Method will give you the method requested by the user like GET, POST etc
	// fmt.Fprint(w, r.Method)

}

//function to handle the root path
func rootHandleFunc(w http.ResponseWriter, r *http.Request) {
	var html, _ = loadFile("welcome.html")
	fmt.Fprint(w, html)
	// fmt.Fprintf(w, "This is the home page")
}

func customerHandlerFunc(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "This is the customer end point")
}
func addHandlerFunc(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w,"This is an add handler function")
}
func main() {
	fmt.Println("Http server")

	//
	http.HandleFunc("/", rootHandleFunc)
	// /user is the endpoint
	http.HandleFunc("/user", userHandlerFunc)

	http.HandleFunc("/customer", customerHandlerFunc)
	//8080  is the port number and nil is the default handler
	//this may give error so we are doing log.Fatal
	log.Fatal(http.ListenAndServe(":8080", nil))
}
