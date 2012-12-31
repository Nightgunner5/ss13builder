package main

import (
	"fmt"
	"os"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "opening file:", err)
		return
	}
	info, err := ReadInfo(f)
	f.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "reading file:", err)
		return
	}
	fmt.Println(info)
}
