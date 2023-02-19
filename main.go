package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
)

func main() {
	if os.Getenv("PPROF") != "" {
		f, err := os.Create(os.Getenv("PPROF") + ".prof")
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
	}

	stdout := bufio.NewWriter(os.Stdout)

	status := run(os.Args, stdout, os.Stderr)
	err := stdout.Flush()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		status = 1
	}

	if os.Getenv("PPROF") != "" {
		pprof.StopCPUProfile()
	}

	os.Exit(status)
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	clr := &color{}
	status := 0
	if len(args) < 2 {
		// FIXME: what if there is 2 arguments, but no positional argument?
		fmt.Fprintf(stderr, "Missing positional argument filename.\nUsage: %s <file1> [file2] ...\n", args[0])
		status = 1
		return status
	}

	// total offset in bytes
	i := 0

	// ascii representation of byte
	ascii := [16]byte{}

	for _, arg := range args[1:] {
		if arg == "-nocolor" {
			clr.disable = true
			continue
		}

		if !clr.disable {
			clr.compute()
		}

		file, err := os.Open(arg)
		if err != nil {
			fmt.Fprintf(stderr, "Failed to open %s: %s\n", arg, err.Error())
			status = 1
		}

		reader := bufio.NewReaderSize(file, 10*1024*1024)

		for {
			b, err := reader.ReadByte()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				fmt.Fprintf(stderr, "Failed to read %s: %s\n", arg, err.Error())
				status = 1
				break
			}

			ascii[i%16] = b

			// print offset
			if i%16 == 0 {
				fmt.Fprintf(stdout, "%08x  ", i)
			}

			// print byte
			fmt.Fprintf(stdout, clr.surround("%02x", b)+" ", b)

			// extra space every 4 bytes
			if (i+1)%4 == 0 {
				fmt.Fprint(stdout, " ")
			}

			// print ascii and break line every 16 bytes
			if (i+1)%16 == 0 {
				fmt.Fprint(stdout, "|")
				printAsciiRow(ascii[:i%16], clr, stdout)
				fmt.Fprintln(stdout, "|")
				ascii = [16]byte{}
			}

			i++
		}
	}

	// print offset
	if i%16 != 0 {

		// compute how many spaces are left for padding ascii part
		left := 16 - i%16
		spaces := 3*left + (left-1)/4 + 1

		fmt.Fprint(stdout, strings.Repeat(" ", spaces))
		fmt.Fprint(stdout, "|")
		printAsciiRow(ascii[:i%16], clr, stdout)
		fmt.Fprintln(stdout, "|")
		fmt.Fprintf(stdout, "%08x\n", i)
	}

	return status
}

func printAsciiRow(ascii []byte, clr *color, stdout io.Writer) {
	var s string
	for _, b := range ascii {

		// is visible ascii
		if b >= 33 && b <= 126 {
			s = clr.surround(string(b), b)
		} else {
			s = clr.surround(".", b)
		}

		fmt.Fprint(stdout, s)
	}
}

type color struct {
	disable bool
	values  [256]string
}

func (c *color) compute() {
	const WHITE_B = "\033[1;37m"

	for i := 0; i < 256; i++ {
		var fg, bg string

		// colors that are very hard to read on a dark background
		barelyVisible := i == 0 || (i >= 16 && i <= 20) || (i >= 232 && i <= 242)

		if barelyVisible {
			fg = WHITE_B + "\033[38;5;" + "255" + "m"
			bg = "\033[48;5;" + strconv.Itoa(int(i)) + "m"

		} else {
			fg = WHITE_B + "\033[38;5;" + strconv.Itoa(int(i)) + "m"
			bg = ""
		}

		c.values[i] = bg + fg
	}
}

func (c *color) surround(s string, clr byte) string {
	const NO_COLOR = "\033[0m"
	return c.values[clr] + s + NO_COLOR
}
