package main

import "fmt"

//here the function getVal returns two int values
func getVal() (int, int) {

	return 3, 4
}
func main() {
	a, b := getVal()
	fmt.Println("The returned values are : ", a, b)

	//if we need a single value only then we can make use of _ to discard the return value
	_, c := getVal()

	fmt.Println(c)

}
