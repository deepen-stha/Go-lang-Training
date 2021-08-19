package main

import (
	"fmt"
	_ "net/http/pprof"
	"test3/packages/Router"
)

//this is the main function which is initiated at the beggining
func main() {

	fmt.Println("Starting the server")
	//calling the function to start our server
	Router.StartServer()
}
