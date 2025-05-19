package main
import "fmt"

func factorial(n int) int {
	if n==0{
		return 1
	}
	return n * factorial(n-1)
}

func main(){
	var n int 
	fmt.Scanf("%d",&n)
	result := factorial(n)
	fmt.Println("Factorial of ", n, "is", result)
}