package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: hx <filepath>\n")
		os.Exit(1)
	}

	filepath := os.Args[1]

	file, err := os.OpenFile(filepath, os.O_RDONLY, 0666) // FIXME: perm
	if err != nil {
		fmt.Printf("Failed to open file %s: %v\n", filepath, err)
		os.Exit(1)
	}

	bytesRead := make([]byte, 1024*1024)
	totalOffset := 0

	for {
		n, err := file.Read(bytesRead)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			fmt.Printf("Failed to read %d bytes from file %s: %v\n", len(bytesRead), filepath, err)
			os.Exit(1)
		}

		prevLine := ""
		printedAsterisk := false

		for currentOffset := 0; currentOffset < n; {
			line, bytesCount := getDumpLine(totalOffset, currentOffset, bytesRead[:n])

			addr := fmt.Sprintf("%s%08x:%s   ", "\033[1m", totalOffset, "\033[0m")

			if line != prevLine {
				fmt.Print(addr, line)
				prevLine = line
				printedAsterisk = false

			} else if !printedAsterisk {
				fmt.Println("*")
				printedAsterisk = true
			}

			currentOffset += bytesCount
			totalOffset += bytesCount
		}
	}

	fmt.Print("\033[0m")
	fmt.Println()
}

func getDumpLine(fileOffset int, start int, bytesRead []byte) (string, int) {

	count := len(bytesRead) - start
	if count > 16 {
		count = 16
	}

	res := ""

	// this is a workaround because I can't use len(res)
	charCount := 0

	ascii := ""

	for i := 0; i < count; i++ {
		b := bytesRead[start+i]
		color := fmt.Sprintf("\033[1;37m\033[38;5;%dm", b)
		if b == 0 || (b >= 16 && b <= 18) || (b >= 232 && b <= 242) {
			color = fmt.Sprintf("\033[1;37m\033[48;5;%dm\033[38;5;255m", b)
		}

		if b >= 32 && b <= 126 {
			ascii += color + string(b) + "\033[0m"
		} else {
			ascii += "."
		}

		res += fmt.Sprintf("%s%02x%s ", color, b, "\033[0m")
		charCount += 3

		if (i+1)%4 == 0 {
			res += " "
			charCount += 1
		}
	}

	maxLength := 53

	res += strings.Repeat(" ", maxLength-charCount)
	res += "\033[0m"
	res += fmt.Sprintf("|%s|\n", ascii)

	return res, count
}
