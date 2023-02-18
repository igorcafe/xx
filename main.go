package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	status := run(os.Args, os.Stdout, os.Stderr)
	os.Exit(status)
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	status := 0
	if len(args) < 2 {
		fmt.Fprintf(stderr, "Missing positional argument filename.\nUsage: %s <file1> [file2] ...\n", args[0])
		status = 1
		return status
	}

	// total offset in bytes
	i := 0

	// ascii representation of byte
	ascii := [16]byte{}

	for _, fname := range args[1:] {
		file, err := os.Open(fname)
		if err != nil {
			fmt.Fprintf(stderr, "Failed to open %s: %s\n", fname, err.Error())
			status = 1
		}

		reader := bufio.NewReader(file)

		for {
			b, err := reader.ReadByte()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				fmt.Fprintf(stderr, "Failed to read %s: %s\n", fname, err.Error())
				status = 1
				break
			}

			if b >= 33 && b <= 126 {
				ascii[i%16] = b
			} else {
				ascii[i%16] = '.'
			}

			// print offset
			if i%16 == 0 {
				fmt.Fprintf(stdout, "%08x  ", i)
			}

			// print byte
			fmt.Fprintf(stdout, "%02x ", b)

			// extra space every 4 bytes
			if (i+1)%4 == 0 {
				fmt.Fprint(stdout, " ")
			}

			// break line every 16 bytes
			if (i+1)%16 == 0 {
				fmt.Fprintln(stdout, "|"+string(ascii[:])+"|")
			}

			i++
		}
	}

	// print offset
	if i%16 != 0 {

		// compute how many spaces are left for padding ascii part
		left := 16 - i%16
		spaces := 3*left + (left-1)/4 + 1

		fmt.Fprintln(stdout, strings.Repeat(" ", spaces)+"|"+string(ascii[:i%16])+"|")
		fmt.Fprintf(stdout, "%08x\n", i)
	}

	return status
}
