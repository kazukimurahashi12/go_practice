package main

import "fmt"

// ジェネリック型エイリアス
type MyAlias[T int | string] = T

func main() {
	fmt.Println("\n--pracice1--")
	practice1()
	fmt.Println("\n--pracice2--")
	practice2()
}

func practice1() {
	var num1 MyAlias[int] = 10
	var num2 MyAlias[int] = 20
	var num3 MyAlias[int] = 30

	var str1 MyAlias[string] = "Go 1.24!"
	var str2 MyAlias[string] = "Generics"
	var str3 MyAlias[string] = "Hello World!"

	fmt.Println("Integers:")
	fmt.Println(num1)
	fmt.Println(num2)
	fmt.Println(num3)

	fmt.Println("\nStrings:")
	fmt.Println(str1)
	fmt.Println(str2)
	fmt.Println(str3)
}

func practice2() {
	// int型スライス
	numbers := []MyAlias[int]{10, 20, 30, 40, 50}

	// string型スライス
	messages := []MyAlias[string]{
		"Go 1.24!",
		"Generics",
		"Let's build something great!",
	}

	fmt.Println("Numbers:")
	for _, num := range numbers {
		fmt.Println(num)
	}

	fmt.Println("\nMessages:")
	for _, msg := range messages {
		fmt.Println(msg)
	}
}
