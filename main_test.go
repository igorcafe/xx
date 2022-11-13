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

func BenchmarkGetDumpLine(b *testing.B) {
	buf := make([]byte, 16)

	for i := 0; i < 16; i++ {
		buf[i] = 0x5a
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = getDumpLine(buf)
	}
}

func BenchmarkDumpBuffer(b *testing.B) {
	buf := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	os.Stdout = nil

	for i := 0; i < b.N; i++ {
		dumpBuffer(0, buf)
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
