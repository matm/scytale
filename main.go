package main

import (
	"fmt"
	"os"
	"strings"
)

var alphabet = []string{
	"A", "B", "C", "D", "E", "F", "G", "H",
	"I", "J", "K", "L", "M", "N", "O", "P",
	"Q", "R", "S", "T", "U", "V", "W", "X",
	"Y", "Z",
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s key plain\n", os.Args[0])
		os.Exit(2)
	}
	vgtable := make([][]string, 26)
	key := strings.ToUpper(os.Args[1])
	clear := strings.ToUpper(os.Args[2])
	fmt.Println("CLEAR", clear)
	key := make([]string, len(clear))
	for k := 0; k < len(clear); k++ {
	}
	if len(key) < len(clear) {
	}
	fmt.Println("KEY  ", key)
	for k := 0; k < 26; k++ {
		vgtable[k] = make([]string, 26)
		for j := 0; j < 26; j++ {
			vgtable[k][j] = alphabet[(j+k)%26]
		}
	}
	fmt.Println("CIPHER", cipher)
}
