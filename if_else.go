package main

import (
	"fmt"
	"os"
	"bufio"
)


func main() {
	var name string

	var age int
	fmt.Println("Enter your name:")
	fmt.Scanln(&name)
	fmt.Println("Enter your age:")
	fmt.Scanln(&age)
	fmt.Println(name, "\n", age)

	if age > 18 {

		fmt.Println("You are an adult")
	} else {
		fmt.Println("You are not adult")
	}

	if 2 > 4 {
		fmt.Println("2 is greater than 4")
	} else if 2 < 4 {
		fmt.Println("2 is less than 4")
	}

	d:=10
	sex:="male"
	if d==10 && sex=="male"{
		fmt.Println("OKKKK")
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter a line")
	line , _ := reader.ReadString('\n')
	fmt.Println("You entered:", line)


}
