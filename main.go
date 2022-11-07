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

	bytesRead := make([]byte, 1024)
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

		for currentOffset := 0; currentOffset < n; {
			printed := printDumpLine(totalOffset, currentOffset, bytesRead[:n])
			currentOffset += printed
			totalOffset += printed
		}
	}

	fmt.Print("\033[0m")
	fmt.Println()
}

func printDumpLine(fileOffset int, start int, bytesRead []byte) int {
	fmt.Printf("%s%08x:%s   ", "\033[1m", fileOffset, "\033[0m")

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

		if b >= 32 && b <= 127 {
			ascii += string(b)
		} else {
			ascii += "."
		}

		high := (b & 0b11110000) >> 4
		low := b & 0b00001111

		highBold := (high / 8) % 2
		lowBold := (low / 8) % 2

		highColor := fmt.Sprintf("\033[%d;%dm", highBold, high%7+31)
		lowColor := fmt.Sprintf("\033[%d;%dm", lowBold, low%7+31)

		res += fmt.Sprintf("%s%x%s%x ", highColor, high, lowColor, low)
		charCount += 3

		if (i+1)%4 == 0 {
			res += " "
			charCount += 1
		}
	}

	maxLength := 53

	fmt.Print(res)
	fmt.Print(strings.Repeat(" ", maxLength-charCount))
	fmt.Print("\033[0m")
	fmt.Printf("|%s|", ascii)
	fmt.Println()

	return count
}
