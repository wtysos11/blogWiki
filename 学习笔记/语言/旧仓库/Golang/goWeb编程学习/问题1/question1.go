package main

import (
    "fmt"
)

type Human struct{
    Name string
    age int
    phone string
}

type Student struct{
    Human
    school string
    loan float32
}

func (h Human) SayHi(){
    fmt.Printf("Hi, I am %s you can call me on %s \n",h.Name,h.phone)
}

func (h Human) Sing(lyrics string){
    fmt.Println("La la la la ...",lyrics)
}

type Men interface{
    SayHi()
    Sing(lyrics string)
}

func main(){
    mike := Student{Human{"Mike", 25, "222-222-XXX"}, "MIT", 0.00}

    var i Men

	i = mike
	fmt.Println(mike.Human.Name)
	fmt.Println(i) //{Human{"Mike", 25, "222-222-XXX"}, "MIT", 0.00}
	mike.Human.Name = "test"
	fmt.Println(i) //{Human{"Mike", 25, "222-222-XXX"}, "MIT", 0.00}
}