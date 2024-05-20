package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
)

const chunkSize = 16

var octetsPerGroup = flag.Int("g", 2, "number of octets per group in normal output. Default 2 (-e: 4)")
var isLittleEndian = flag.Bool("e", false, "little-endian dump")
var outputLength = flag.Int("l", -1, "stop after <l> octets")

var padding int

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		os.Exit(1)
	}
	if *octetsPerGroup < 0 {
		*octetsPerGroup = 4
	} else if *octetsPerGroup == 0 || *octetsPerGroup > 16 {
		*octetsPerGroup = 16
	}
	if *isLittleEndian && !isPowerOfTwo(*octetsPerGroup) {
		fmt.Println("xxd: number of octets per group must be a power of 2 with -e.")
		os.Exit(1)
	}
	if *outputLength < 0 {
		*outputLength = math.MaxInt
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
	var bytesRead int
	var lengthReached bool
	for i := 0; !lengthReached; i++ {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if err == io.EOF {
			break
		}

		bytesLine := buf[:n]
		if n >= *outputLength-bytesRead {
			bytesLine = buf[:*outputLength-bytesRead]
			lengthReached = true
		}
		bytesRead += n
		fmt.Printf("%08x: ", i*chunkSize)
		padding = 32 + 16 / *octetsPerGroup

		if 16%*octetsPerGroup == 0 {
			padding--
		}

		printHex(bytesLine)
		fmt.Printf("%s  ", strings.Repeat(" ", padding))
		printText(bytesLine)
		fmt.Println()
	}
}

func isPowerOfTwo(n int) bool {
	if n > 0 && n&(n-1) == 0 {
		return true
	}
	return false
}

func printHex(bytes []byte) {
	if *isLittleEndian {
		groupCount := int(math.Ceil(float64(len(bytes)) / float64(*octetsPerGroup)))
		for i := range groupCount {
			for j := *octetsPerGroup - 1; j >= 0; j-- {
				if j+i**octetsPerGroup >= len(bytes) {
					fmt.Print("  ")
					padding -= 2
					continue
				}
				fmt.Printf("%02x", bytes[j+i**octetsPerGroup])
				padding -= 2
			}
			if i != groupCount-1 {
				fmt.Print(" ")
				padding -= 1
			}
		}
	} else {
		for i, b := range bytes {
			fmt.Printf("%02x", b)
			padding -= 2
			if (i+1)%*octetsPerGroup == 0 && i != len(bytes)-1 {
				fmt.Print(" ")
				padding -= 1
			}
		}
	}
}

func printText(bytes []byte) {
	for _, b := range bytes {
		if string(b) == "\n" || b == 0 {
			fmt.Print(".")
			continue
		}
		fmt.Print(string(b))
	}
}
