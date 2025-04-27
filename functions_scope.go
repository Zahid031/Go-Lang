package main

import "fmt"

var (
	a=10
	b=20
)

func add(x int, y int) {
	res := x+y
	PrintSum(res)
}

func PrintSum(z int){
	fmt.Println("Sum: ", z)
}

func main() {
	add(a,b)
}