package main
import "fmt"
func main(){
	arr:=[...]int{1,2,3,4,5,6,7,8,9,10}
	for i:=0; i<len(arr); i++{
		fmt.Println(arr[i])
	}
	for _,value := range arr {
		//fmt.Println("Index:", idx, "Value:", value)
		fmt.Println( "Value:", value)

	}
}