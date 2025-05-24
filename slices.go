package main
import "fmt"

func main(){
	var slice1=[]int{1,2,3,4,5}//slice literal
	fmt.Println(slice1)
	fmt.Println(len(slice1),cap(slice1))

	//creating slice from an array
	arr:=[...]int{1,2,3,4,5}
	//creating slice with make
	slice3:=make([]int,5,5)
	slice2:=arr[:4]
	fmt.Println(slice2)
	fmt.Println(slice3)
	fmt.Println((len(slice3)),"capacity:",cap(slice3))
	fmt.Printf("%d\n",len(slice3))
	slice3=append(slice3, slice1...)
	slice3=append(slice3, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	fmt.Println((slice3))
	numbers := []int{1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	fmt.Printf("numbers: %v\n", numbers)
	fmt.Printf("Capacity: %d\n", cap(numbers))
	fmt.Printf("Length: %d\n", len(numbers))
	needNumbers := numbers[:len(numbers)-10]
	fmt.Printf("needNumbers: %v\n", needNumbers)
	numbersCopy := make([]int,len(needNumbers))
	copy(numbersCopy, needNumbers)
	fmt.Printf("numbersCopy: %v\n", numbersCopy)
	ar:=[...]int{1,2,3,4,5,6,7,8,9,10}
	s1:=ar[1:6]
	s2:=s1[2:4]
	fmt.Println(s1,len(s1),cap(s1))
	fmt.Println(s2,len(s2),cap(s2))

	var s []int //nil slice or empty slice
	s = append(s, 1)//if s nil then it will create a array
	s = append(s, 2)
	fmt.Println(s)



}