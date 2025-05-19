package main
import "fmt"

func add(a int, b int) int{
	return a+b
}

func getvalues(a , b int)(int, int){
	addition := a+b
	div := a*b
	return addition, div
}

func welcomemessage(){
	fmt.Println("Welcome to Golang")
}

func getUserName() string{
	var name string
	fmt.Println("Enter your name: ")
	fmt.Scanln(&name)
	return name
}

func getTwoNumbers() (int, int){
	var a,b int
	fmt.Println("Enter first number: ")
	fmt.Scanln(&a)
	fmt.Println("Enter second number: ")
	fmt.Scanln(&b)
	return a,b
}

func display(name string, sum int){
	fmt.Println("Hello: ",name)
	fmt.Println("Sume of two number is: ",sum)
}



func main(){
	welcomemessage()
	name := getUserName()
	num1 , num2 := getTwoNumbers()
	sum:=add(num1,num2)
	display(name, sum)

	func(a int, b int){
		c :=a+b
		fmt.Println("Sum of two..:", c)
	}(3,4)
}


func init(){
	fmt.Println("This is init function")
}


