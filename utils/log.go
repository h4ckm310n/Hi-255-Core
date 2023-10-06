package utils

import "fmt"

func Log(content ...interface{}) {
	fmt.Println(content)
}

func Err(content ...interface{}) {
	fmt.Print("Error: ")
	fmt.Println(content)
}
