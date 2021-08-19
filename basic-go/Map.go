package main

import (
	"fmt"
)

func main() {
	mapit()
}

func mapit() {

	//making a new map with key- type as string and value type as string
	var studentlist = make(map[string]string)

	studentlist["name"] = "Deepen"
	studentlist["age"] = "21"
	studentlist["department"] = "CSE"

	fmt.Println("before deleting any key: ", studentlist)

	delete(studentlist, "age")

	fmt.Println("after deleting key : ", studentlist)

}
