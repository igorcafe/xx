package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const WHITE_B = "\033[1;37m"
const NO_COLOR = "\033[0m"
const NO_COLOR_B = "\033[1m"

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

	bytesRead := make([]byte, 1024*1024) // TODO: optimal buffer size
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

		dumpBuffer(totalOffset, bytesRead[:n])
	}

	fmt.Print(NO_COLOR)
	fmt.Println()
}

func dumpBuffer(totalOffset int, bytesRead []byte) {
	prevLine := ""
	printedAsterisk := false

	for currentOffset := 0; currentOffset < len(bytesRead); {
		line, bytesCount := getDumpLine(bytesRead[currentOffset:])

		addr := fmt.Sprintf("%s%08x:%s   ", NO_COLOR_B, totalOffset, NO_COLOR)

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

func getDumpLine(bytesRead []byte) (string, int) {
	count := len(bytesRead)
	if count > 16 {
		count = 16
	}

	resBuilder := strings.Builder{}
	resBuilder.Grow(413) // FIXME: magic number

	// this is a workaround because I can't use len(res)
	charCount := 0

	asciiBuilder := strings.Builder{}
	asciiBuilder.Grow(count)

	for i := 0; i < count; i++ {
		b := bytesRead[i]
		fg := color256(b, true)
		bg := ""

		// colors that are barely readable with dark background
		useBrightBg := b == 0 || (b >= 16 && b <= 18) || (b >= 232 && b <= 242)

		if useBrightBg {
			fg = color256(255, true)
			bg = color256(b, false)
		}

		color := WHITE_B + bg + fg

		isPrintableAscii := b >= 32 && b <= 126

		if isPrintableAscii {
			asciiBuilder.WriteString(color + string(b) + NO_COLOR)
		} else {
			asciiBuilder.WriteString(".")
		}

		resBuilder.WriteString(color)
		resBuilder.WriteString(fmt.Sprintf("%02x", b)) // TODO: func with 0 alloc/op
		resBuilder.WriteString(NO_COLOR)
		resBuilder.WriteString(" ")
		charCount += 3

		if (i+1)%4 == 0 {
			resBuilder.WriteString(" ")
			charCount += 1
		}
	}

	maxLength := 53 // FIXME: magic number

	resBuilder.WriteString(strings.Repeat(" ", maxLength-charCount))
	resBuilder.WriteString(NO_COLOR)
	resBuilder.WriteString(fmt.Sprintf("|%s|\n", asciiBuilder.String()))

	return resBuilder.String(), count
}

func color256(b byte, foreground bool) string {
	colorType := 48
	if foreground {
		colorType = 38
	}

	return fmt.Sprintf("\033[%d;5;%dm", colorType, b)
}
