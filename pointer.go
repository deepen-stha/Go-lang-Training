package main

import "fmt"

func demo(val *int) {
	*val = 20
	fmt.Println(*val)
}

func main() {
	var i int = 5
	b := &i
	fmt.Println(b)

	//incrementing the value by making use of pointer
	*b++
	fmt.Println(i)

	//passisng a pointer to a function
	demo(b)
	fmt.Println(*b) //now the value of the b is changed to 20

}
