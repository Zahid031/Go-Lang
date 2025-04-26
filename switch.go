package main
import "fmt"

func main(){
	var a int
	fmt.Scanf("%d",&a)
	switch a {
	case 1:
		fmt.Println("Entered 1")
	case 2:
		fmt.Println("Entered 2")
	default:
		fmt.Println("Entered ....")
	}
}