package main

import (
	"fmt"

	"github.com/jstang9527/gateway/thirdpart/selm2"
)

func main() {
	path := []string{"div[1]", "a[1]", "div[1]", "li[3]", "ul[1]", "div[2]", "div[1]", "div[1]"}
	link := []string{"a[1]", "div[1]", "li[3]", "ul[1]", "div[2]", "div[1]", "div[1]"}
	fmt.Println("path: ", path)
	fmt.Println("link: ", link)
	a1 := selm2.Reverse(path)
	a2 := selm2.Reverse(link)
	fmt.Println(a1)
	fmt.Println(a2)
}
