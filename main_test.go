package main

import (
	"bufio"
	"strings"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func BenchmarkParser(b *testing.B) {
	delta := 10000
	x := "add 1 1\nprint 1\n"
	var l = delta
	s := strings.Repeat(x, delta)

	println()
	for i := 1; i <= 40; i++ {
		eventLoop := new(EventLoop)
		eventLoop.Start()
		s += strings.Repeat(x, delta)

		scanner := bufio.NewScanner(strings.NewReader(s))

		start := time.Now()
		eventLoop.parseScanner(scanner)
		end := time.Now()
		println(l, end.Sub(start)/time.Millisecond)

		delta = int(float32(l) * 0.2)
		l = int(float32(l) * 1.2)
	}
}
