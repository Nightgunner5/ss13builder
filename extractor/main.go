package main

import (
	"fmt"
	"os"
	"strings"
	"image/png"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "opening dmi:", err)
		return
	}
	info, err := ReadInfo(f)
	f.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "reading dmi:", err)
		return
	}

	f, err = os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "opening dmi:", err)
		return
	}
	img, err := png.Decode(f)
	f.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "decoding png:", err)
		return
	}

	dir := os.Args[1] + ".dir"
	os.Mkdir(dir, 0755)
	Extract(dir, strings.NewReader(info), img)
}
