//package clause start every source file. main is a special name declaring an executable rather than a lib
package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
	beyondHello()
}

func beyondHello() {
	var x int
	x = 3
	y := 4	// short declaration -> infer the type
	sum, prod := learnMultiple(x, y) // func return 2 values
	fmt.Println("sum:", sum, "prod:", prod)
	learnTypes()
}

func learnMultiple(x, y int) (sum, prod int) { // return value or signature
	return x + y, x * y
}

func learnTypes() {
	str := "Learn Go!"
	s2 := `A "raw" string literal
	can include line breaks`

	n := byte('\n') // conversion syntax with short declaration : byte is an alias for uint8
	var a4 [4]int // array of 4 int, init to all 0
	a5 := [...]int{3, 1, 5, 10, 100} // array init with fixed size of 5 elements
	// array have value semantics
	a4_cpy := a4 // a4_cpy is a copy of a4, 2 separate instances
	a4_cpy[0] = 25 // only copy is changed
	fmt.Println(a4_cpy[0] == a4[0]) // False
	// slices have dynamic size
	s3 := []int{4, 5, 9} // compare to a5. No ellipsis here.
	s4 := make([]int, 4) // allocates slice of 4 ints, init at 0
	var d2 [][]float64 // declaration only, nothing allocated here
	bs := []byte("a slice") // type conversion syntax
	// slice (as well as maps and channels) have references semantics
	s3_cpy := s3 // both var point to the same instance
	s3_cpy[0] = 0 // both are updated
	fmt.Println(s3_copy[0] == s3[0])

	// because they are dynamic, slices can be appended on demand, using the builin append()
	s := []int{1, 2, 3}
	s = append(s, 4, 5, 6)
	fmt.Println(s)
	// we can also, instead of a list of atomic element pass a slice reference, unpacking it with trailing ellipsis
	s = append(s, []int{7, 8, 9}...)
	fmt.Println(s)

	p, q := learnMemory() // declares p, q to be type pointer to int
	fmt.Println(*p, *q) // * follow a pointer

	// maps are a dynamically growable associative array type


}

func learnMemory() (p, q *int) {
	p = new(int) // builtin to allocate new memory
	s := make([]int, 20) // allocate 20 ints as a single block of memory
	s[3] = 7
	r := -2
	return &s[3], &r // & -> takes the address of an object
}
