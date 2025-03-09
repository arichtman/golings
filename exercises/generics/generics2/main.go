// generics2
// Make me compile!

package main

import "fmt"

type Number interface {
	int | float64
}

// Can't fulfil interface with method definition on nonlocal type
// func (f float64) int() {
// 	int(f)
// }

func main() {
	fmt.Println(addNumbers(1, 2))
	fmt.Println(addNumbers(1.0, 2.3))
}

func addNumbers[T Number](n1 T, n2 T) T {
	return n1 + n2
}
