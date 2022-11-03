package logutil

import (
	"bufio"
	"bytes"
	"io"
)

func Forward(printf func(fmt string, args ...any)) io.Writer {
	return &outputForward{
		printf: printf,
	}
}

type outputForward struct {
	printf func(fmt string, args ...any)
}

func (o *outputForward) Write(p []byte) (n int, err error) {
	s := bufio.NewScanner(bytes.NewBuffer(p))
	for s.Scan() {
		if line := s.Text(); len(line) > 0 {
			o.printf("%s", line)
		}
	}
	return len(p), nil
}
