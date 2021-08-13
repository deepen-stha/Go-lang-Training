package main

import (
	"fmt"
)

func main() {
	fmt.Println("Enter the element in an slice:")
	var n int

	fmt.Scanln(&n)

	s := make([]int, n)

	fmt.Println("Enter the values in slice:")

	for i := 0; i < len(s); i++ {
		fmt.Scanln(&s[i])
	}

	//creating new slice variable to store the reversed slice

	newslice := make([]int, 0)

	fmt.Println(len(s))
	for i := len(s) - 1; i >= 0; i-- {

		newslice = append(newslice, s[i])
	}
	fmt.Println("Before reversing slice was :", s)
	fmt.Println("after reversing :", newslice)

}
