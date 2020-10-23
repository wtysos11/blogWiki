package main

import "reflect"
import "fmt"
func main() {
    t := reflect.TypeOf(make(map [string]string)).Elem()
    fmt.Println(t)
}