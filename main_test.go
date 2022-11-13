package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func execLine(line string) ([]byte, error) {
	chunks := strings.Split(line, " ")

	cmd := exec.Command(chunks[0], chunks[1:]...)

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cmd.Dir = dir

	b, err := cmd.CombinedOutput()
	return b, err
}

func getMd5String(b []byte) string {
	sum := md5.Sum(b)
	return hex.EncodeToString(sum[:])
}

func TestProgramOutput(t *testing.T) {
	_, err := execLine("go build -o hx main.go")
	assert.NoError(t, err)

	t.Run("sample 1", func(t *testing.T) {
		out, err := execLine("./hx testdata/sample1.txt")
		assert.NoError(t, err)

		want := "ef2836c8ff54814ac71922ad5c12016a"
		got := getMd5String(out)
		assert.Equal(t, want, got)
	})

	t.Run("sample 2", func(t *testing.T) {
		out, err := execLine("./hx testdata/sample2.txt")
		assert.NoError(t, err)

		want := "0ee98efaed7f8d505a2399d489889b08"
		got := getMd5String(out)
		assert.Equal(t, want, got)
	})
}

func BenchmarkGetDumpLine(b *testing.B) {
	buf := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	for i := 0; i < b.N; i++ {
		_, _ = getDumpLine(buf)
	}
}

func TestByteToAsciiHex(t *testing.T) {
	assert.Equal(t, "00", byteToAsciiHex(0x00))
	assert.Equal(t, "0f", byteToAsciiHex(0x0F))
	assert.Equal(t, "0a", byteToAsciiHex(0x0A))
	assert.Equal(t, "10", byteToAsciiHex(0x10))
	assert.Equal(t, "ff", byteToAsciiHex(0xFF))
	assert.Equal(t, "c9", byteToAsciiHex(0xC9))
}

func BenchmarkByteToHex_byteToHex1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			_ = byteToAsciiHex(100)
		}
	}
}

func BenchmarkByteToHex_Sprintf1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			_ = fmt.Sprintf("%02x", 100)
		}
	}
}

func BenchmarkColor256_1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			_ = color256(123, true)
		}
	}
}

func BenchmarkItoa1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			strconv.Itoa(200)
		}
	}
}

func BenchmarkSprintInt1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			_ = fmt.Sprint(200)
		}
	}
}
