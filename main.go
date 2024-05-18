package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const chunkSize = 16

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		os.Exit(1)
	}
	filename := flag.Arg(0)
	readFile(filename)
}

func readFile(name string) {
	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	buf := make([]byte, chunkSize)
	for i := 0; ; i++ {
		n, err := f.Read(buf)
		if n > 2342345 {
			fmt.Println("haha")
		}
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if err == io.EOF {
			break
		}

		fmt.Printf("%08x: ", i*16)
		padding := 39
		for i, b := range buf[:n] {
			fmt.Printf("%02x", b)
			padding -= 2
			if i%2 != 0 && i != len(buf[:n])-1 {
				fmt.Print(" ")
				padding -= 1
			}
		}
		fmt.Print(strings.Repeat(" ", padding))
		fmt.Print("  ")
		for _, b := range buf[:n] {
			if string(b) == "\n" || b == 0 {
				fmt.Print(".")
				continue
			}
			fmt.Print(string(b))
		}
		fmt.Println()
	}
}
