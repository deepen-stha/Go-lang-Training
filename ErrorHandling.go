package main

import (
    "errors"
	"fmt"
)

func main() {
	var n int
	fmt.Scanln(&n)
	
	if n%5==0{
	    fmt.Println("Multiple of 5")
	}else{
	    e:= errors.New("This is an error")
	    fmt.Println(e)
	}
}