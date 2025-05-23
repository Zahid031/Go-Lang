package main

import "fmt"

type User struct {
	Name string
	Age int
}

func printDetails(usr User){
	fmt.Println("Name: ",usr.Name)
	fmt.Println("Age: ",usr.Age)
}

//Receiver function

func (usr User) printDetails(){
	fmt.Println("Name: ",usr.Name)
	fmt.Println("Age: ",usr.Age)
}

func (user User) updateName(name string){
	user.Name = name
	fmt.Println("Name: ",user.Name)
	fmt.Println("Age: ",user.Age)
}
func main(){

	user :=User{
		Name : "Zahid",
		Age: 25,
	}
	printDetails(user)
	user.printDetails()
	user.updateName("Nahid")
}