package main

import (
	"fmt"
)

func main() {
	fmt.Println("Enter the value of n:")
	var n int

	fmt.Scanln(&n)

	const temp1 string = "Fizz"
	const temp2 string = "Buzz"
	for i := 1; i <= n; i++ {
		if i%3 == 0 && i%5 == 0 {
			fmt.Println(temp1 + temp2)
		} else if i%5 == 0 {
			fmt.Println(temp2)
		} else if i%3 == 0 {
			fmt.Println(temp1)
		} else {
			fmt.Println(i)
		}

	}
}