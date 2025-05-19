package main
import "fmt"
func main(){
	var arr = [5]int{1,2,3,4,5}
	var arr2 = [...]string{"RTR", "BMW", "TESLA"}
	//initializing specific element in an array
	arr3:=[...]int{1:10, 8:5}
	fmt.Println(arr)
	fmt.Println(arr2[2])
	fmt.Println(arr3)
	fmt.Println(len(arr3))
}