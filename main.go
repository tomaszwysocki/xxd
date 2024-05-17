package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		os.Exit(1)
	}
	filename := flag.Arg(0)
	readFile(filename)
}

func readFile(name string) {
	data, err := os.ReadFile(name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i := 0; i < len(data); i += 16 {
		var padding string
		end := i + 16

		if end > len(data) {
			end = len(data)
		}

		bytes := data[i:end]
		fmt.Printf("%08x: ", i)

		var fourBytesSlice []string

		for j := 0; j < end-i; j += 2 {
			endInner := j + 2
			if endInner > end-i {
				endInner = end - i
			}
			fourBytesSlice = append(fourBytesSlice, fmt.Sprintf("%x", bytes[j:endInner]))
		}

		fmt.Printf("%v", strings.Trim(fmt.Sprint(fourBytesSlice), "[]"))

		if len(fourBytesSlice) < 8 {
			padding = fmt.Sprintf(strings.Repeat(" ", (16-(end-i))*2+8-len(fourBytesSlice)))
		}
		fmt.Printf("%v", padding)

		var bytesString string
		for _, b := range bytes {
			if string(b) == "\n" || b == 0 {
				bytesString += fmt.Sprint(".")
				continue
			}
			bytesString += fmt.Sprint(string(b))
		}
		fmt.Printf("  %v", bytesString)
		fmt.Println()
	}
}
