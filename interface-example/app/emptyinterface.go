package main

import "fmt"

func main() {
	// declaring an interface
	// ----
	type Empty interface {
	}

	// interface variable
	// ----
	var e Empty
	describe(e)

	e = 7
	describe(e)

	e = "Hello, World!"
	describe(e)

	e = func() int {
		return 1
	}
	describe(e)

	var i interface{}
	describe(i)

	i = 42
	describe(i)

	i = "hello"
	describe(i)

	PrintBasedOnType(10)
	PrintBasedOnType("I'm a string")
}

func describe(i interface{}) {
	fmt.Printf("e's value: %v, type: %T\n", i, i)
}

func PrintBasedOnType(i interface{}) {
	// A type switch can work out what type an interface represents.
	switch x := i.(type) {
	case int:
		fmt.Printf("i is an integer: %d\n", x) // Print the integer.
	case string:
		fmt.Printf("i is a string: %s\n", x) // Print the string.
	}
}

// https://golangbot.com/interfaces-part-1/