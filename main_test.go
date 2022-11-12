package main

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"os/exec"
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