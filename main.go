package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const WHITE_B = "\033[1;37m"
const NO_COLOR = "\033[0m"
const NO_COLOR_B = "\033[1m"

const ASCII_0_POSITION = byte(48)
const ASCII_A_POSITION = byte(97)

var colors [256]string

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: hx <filepath>\n")
		os.Exit(1)
	}

	precomputeColors()
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

		addr := NO_COLOR_B + fmt.Sprintf("%08x:", totalOffset) + NO_COLOR + "   "

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
	resBuilder.Grow(208) // FIXME: magic number

	// this is a workaround because I can't use resBuilder.Len()
	charCount := 0

	for i := 0; i < count; i++ {
		b := bytesRead[i]
		color := colors[b]

		resBuilder.WriteString(color)
		resBuilder.WriteString(byteToAsciiHex(b))
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
	resBuilder.WriteString(NO_COLOR + "|")

	for i := 0; i < count; i++ {
		b := bytesRead[i]
		color := colors[b]

		isPrintableAscii := b >= 32 && b <= 126

		if isPrintableAscii {
			resBuilder.WriteString(color + string(b) + NO_COLOR)
		} else {
			resBuilder.WriteString(".")
		}
	}

	resBuilder.WriteString(NO_COLOR + "|\n")

	l := resBuilder.Len()
	_ = l

	return resBuilder.String(), count
}

func precomputeColors() {
	for i := 0; i < 256; i++ {
		var fg, bg string

		// colors that are very dark
		barelyVisible := i == 0 || (i >= 16 && i <= 20) || (i >= 232 && i <= 242)

		if barelyVisible {
			fg = WHITE_B + "\033[38;5;" + "255" + "m"
			bg = "\033[48;5;" + strconv.Itoa(int(i)) + "m"

		} else {
			fg = WHITE_B + "\033[38;5;" + strconv.Itoa(int(i)) + "m"
			bg = ""
		}

		colors[i] = bg + fg
	}
}

func halfByteToAsciiHex(b byte) string {
	half := b & 0x0F

	if half <= 0x9 {
		return string(half + ASCII_0_POSITION)
	} else {
		return string(half + ASCII_A_POSITION - 0x0a)
	}
}

func byteToAsciiHex(b byte) string {
	low := b & 0x0F
	high := (b & 0xF0) >> 4

	return halfByteToAsciiHex(high) + halfByteToAsciiHex(low)
}
