// generics1
// Make me compile!

package main

import "fmt"

type MyType struct{}

func main() {
	print("Hello, World!")
	print(42)
	print(MyType{})
}

func print(value any) {
	fmt.Println(value)
}
