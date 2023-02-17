package main

import "fmt"

func main() {
	var name string

	fmt.Print("i am rais")
	fmt.Scanln(&name)

	fmt.Printf("Halo, %s!\n", name)
}
