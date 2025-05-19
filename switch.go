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

	switch a {
	case 1,2,3:
		fmt.Println("Entered 1,2,3")
	case 4,5,6:
		fmt.Println("Entered 4,5,6")
	default:
		fmt.Println("Entered more than 6")
	}
}