// slices1
// Make me compile!

package main

import "fmt"

func main() {
	a := make([]int, 3, 10) // play with length and capacity
	// var b [2]int
	// b = [...]int{1, 0}
	fmt.Println("length of 'a':", len(a))
	fmt.Println("capacity of 'a':", cap(a))
}
