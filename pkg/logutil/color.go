package logutil

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"
)

var NoColor = noColorExists()

func noColorExists() bool {
	_, exists := os.LookupEnv("NO_COLOR")
	return exists
}

type WrapWriter = func(w io.Writer) io.Writer

func WithColor(attrs ...Attribute) func(w io.Writer) io.Writer {
	return func(w io.Writer) io.Writer {
		if NoColor {
			return w
		}
		return &colorWriter{w: w, attrs: attrs}
	}
}

type colorWriter struct {
	w     io.Writer
	attrs []Attribute
}

// write SGR sequence "\x1b[...m"
func writeSequenceTo(w io.Writer, attrs ...Attribute) {
	if len(attrs) > 0 {
		_, _ = io.WriteString(w, escape)
		_, _ = io.WriteString(w, "[")
		for i, attr := range attrs {
			if i > 0 {
				_, _ = io.WriteString(w, ";")
			}
			_, _ = io.WriteString(w, strconv.Itoa(int(attr)))
		}
		_, _ = io.WriteString(w, "m")
	}
}

func (c *colorWriter) sequence() string {
	format := make([]string, len(c.attrs))
	for i, v := range c.attrs {
		format[i] = strconv.Itoa(int(v))
	}

	return strings.Join(format, ";")
}

func (c *colorWriter) Write(p []byte) (n int, err error) {
	b := bytes.NewBuffer(nil)

	writeSequenceTo(b, c.attrs...)

	_, _ = b.Write(p)

	writeSequenceTo(b, Reset)

	i, err := io.Copy(c.w, b)
	return int(i), err
}

type Attribute int

const escape = "\x1b"

// Base attributes
const (
	Reset Attribute = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

// Foreground text colors
const (
	FgBlack Attribute = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Foreground Hi-Intensity text colors
const (
	FgHiBlack Attribute = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

// Background text colors
const (
	BgBlack Attribute = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Background Hi-Intensity text colors
const (
	BgHiBlack Attribute = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)
