package main
import "fmt"

var a = 10

func main(){

	if true {
		a:=20 // this is a variable shadowing
		fmt.Println("Inside if block: ", a)
	}
	fmt.Println("Outside block: ", a)

}