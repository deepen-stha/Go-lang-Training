package main

import (
	"fmt"
)

//fact function accepts an integer as argument and return an int
func fact(n int) int {
	if n == 1 {
		return 1
	}
	return n * fact(n-1)

}
func main() {
	var n int
	fmt.Print("Enter the value of n: ")
	fmt.Scanln(&n)
	result := fact(n)
	fmt.Printf("Factorial of %d is: %d\n", n, result)
}
