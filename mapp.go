package main
import "fmt"

func main(){
	var mp = make(map[string]int)
	mp["one"]=1
	mp["two"]= 2
	mp["three"]=3

	fmt.Println(mp)
	delete(mp, "one")
	fmt.Println(mp)
}