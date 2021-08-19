package main

import "fmt"

type Employee struct {
	name   string
	salary int
	age    int
}

func main() {

	e1 := Employee{
		name:   "Deepen",
		salary: 200000000,
		age:    21,
	}

	fmt.Println(e1.age)
	fmt.Println(e1.salary)
	fmt.Println(e1.name)
}
