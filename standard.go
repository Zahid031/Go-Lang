// This is standard function or named function

package main
import "fmt"
func addition(a int,b int) int{
	return a+b
}

func add(){
	c:=2
	d:=3
	fmt.Println((addition(c,d)))

}


func main(){
	add()
}