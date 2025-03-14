package main

import (
	"fmt"
	"go-fundamental/function"
	"go-fundamental/multithreading"
)

func main() {
	var a int = 1
	fmt.Println("Add 2 to a: ", function.Add(a, 2))
	fmt.Println("Multiply 3 numbers: ", function.Multiply(7, 2, 7))
	fmt.Println("Multiply 4 numbers: ", function.Multiply(2, 3, 5, 7))
	fmt.Println("Multiply 5 numbers: ", function.Multiply(6, 9, 4, 2, 0))
	val, err := function.Divide(5, 0)
	fmt.Println("Divide by zero: ", val, " ", err)
	val, err = function.Divide(10, 3)
	fmt.Println("Normal division: ", val, " ", err)
	multithreading.Test()
}
