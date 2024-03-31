package main

import (
	"fmt"

	"example.com/greetings"
)

func main() {
	message := greetings.Hello("Moses")
	fmt.Println(message)
}
