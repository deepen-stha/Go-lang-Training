package main

import (
	"fmt"
)

func main() {
	fmt.Println("Enter the element in an array:")
	var n int

	fmt.Scanln(&n)
	fmt.Println(n)
	// 	n := 10

	var arr [10]int

	var sum int = 0

	fmt.Println("Enter the value in array:")

	for i := 0; i < n; i++ {
		fmt.Scanln(&arr[i])
	}

	for i := 0; i < n; i++ {
		sum = sum + arr[i]
		// 		fmt.Println(i)
	}
	fmt.Println("The sum of element :", arr)
	fmt.Println(sum)
}
