package main
import "fmt"

func print(numbers *[]int){
	fmt.Println("Numbers: ", *numbers)
}

func main(){

	x:=10
	ptr:=&x
	fmt.Println("Value of x: ",x)
	fmt.Println("Address of x: ",ptr)
	fmt.Println("Value of x using pointer: ", *ptr)
	*ptr=20
	fmt.Println("Value of x after using pointer: ",x)
	numbers:=[]int{1,2,3,4,5}
	print(&numbers)
}