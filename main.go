package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "Usage:\n\t%s FILE1 [FILE2...]", os.Args[0])
		os.Exit(1)
		return
	}

	loadColors()

	for _, path := range os.Args[1:] {
		func() {
			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			dump(os.Stdout, file)
		}()
	}
}

func dump(w io.Writer, r io.Reader) error {
	bufR := bufio.NewReader(r)
	bufW := bufio.NewWriter(w)
	offset := int64(-1)
	row := make([]byte, 16)

	defer bufW.Flush()

	for {
		b, err := bufR.ReadByte()
		if errors.Is(err, io.EOF) {
			return dumpRow(bufW, offset, row)
		}
		if err != nil {
			return err
		}

		offset++
		rowOffset := offset % 16
		row[rowOffset] = b

		if offset > 0 && rowOffset == 15 {
			if err := dumpRow(bufW, offset, row); err != nil {
				return err
			}
		}
	}
}

func dumpRow(w io.Writer, offset int64, row []byte) error {
	rowLen := int(offset%16 + 1)

	s := fmt.Sprintf("%08x", offset-int64(rowLen)+1) + "   "
	for i := 0; i < 16; i++ {
		if i < rowLen {
			// args = append(args, row[i])
			s += coloredByte(row[i]) + " "
		} else {
			s += "   "
		}
		if (i+1)%4 == 0 {
			s += " "
		}
	}

	s += " |"
	for i := 0; i < 16; i++ {
		if i < rowLen {
			b := byte('.')
			if row[i] >= 33 && row[i] <= 126 {
				b = row[i]
			}
			s += coloredString(string(b), row[i])
		} else {
			s += " "
		}
	}
	s += "|\n"

	_, err := fmt.Fprint(w, s)

	return err
}

func coloredString(s string, b byte) string {
	// if b == 25 {
	// 	fmt.Println("\n\nSHIT", strings.ReplaceAll(colors[b], "\\", ">"), s)
	// }
	return colors[b] + s + "\033[0m"
}

func coloredByte(b byte) string {
	return fmt.Sprintf("%s%02x%s", colors[b], b, "\033[0m")
}

var colors [256]string

func loadColors() {
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

		colors[i] = bg + fg
	}
}
