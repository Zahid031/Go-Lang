package main
import "fmt"
type User struct {  //member variable or properties Name and Age
	Name string
	Age int
}

func main(){
	var user1 User
	user1.Name="Zahid"
	user1.Age=25
	fmt.Println(user1)
	user2 := User{      //Instantiate a struct // Instance
		Name:"Nahid",
		Age: 30,
	}
	fmt.Println("Name: ",user2.Name)
	fmt.Println("Age: ",user2.Age)
}