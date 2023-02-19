package main

import (
	"bufio"
	"os"
	"testing"
)

func Benchmark_RunWithUnbufferedOutputs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		run([]string{"./xx", "samples/10M"}, os.Stdout, os.Stderr)
	}
}

func Benchmark_RunWithBufferedOutputs(b *testing.B) {
	stdout := bufio.NewWriter(os.Stdout)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		run([]string{"./xx", "samples/10M"}, stdout, os.Stderr)
		stdout.Flush()
	}
}
