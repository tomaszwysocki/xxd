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

var octetsPerGroup = flag.Int("g", 2, "number of octets per group in normal output. Default 2 (-e: 4)")

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
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if err == io.EOF {
			break
		}

		bytesLine := buf[:n]

		fmt.Printf("%08x: ", i*16)
		padding := 32 + 16 / *octetsPerGroup

		if 16%*octetsPerGroup == 0 {
			padding--
		}

		for i, b := range bytesLine {
			fmt.Printf("%02x", b)
			padding -= 2
			if (i+1)%*octetsPerGroup == 0 && i != len(bytesLine)-1 {
				fmt.Print(" ")
				padding -= 1
			}
		}
		fmt.Print(strings.Repeat(" ", padding))
		fmt.Print("  ")
		for _, b := range bytesLine {
			if string(b) == "\n" || b == 0 {
				fmt.Print(".")
				continue
			}
			fmt.Print(string(b))
		}
		fmt.Println()
	}
}
