package logger

import (
	"io"
	"os"
	"sync"

	"github.com/fatih/color"
)

// DefaultWriter var.
var DefaultWriter = &Writer{
	w: os.Stderr,
	c: true,
}

// NewWriter ...
func NewWriter(w io.Writer, supportANSIColors bool) *Writer {
	return &Writer{
		w: w,
		c: supportANSIColors,
	}
}

// NewWriterFile ...
func NewWriterFile(path string) (*Writer, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return nil, err
	}
	return NewWriter(f, false), nil
}

// Writer struct.
type Writer struct {
	w io.Writer
	b []byte
	m sync.Mutex
	c bool
}

func (w *Writer) clone() *Writer {
	w.m.Lock()
	defer w.m.Unlock()
	return &Writer{
		w: w.w,
		c: w.c,
	}
}

var (
	white  = color.New(color.FgHiWhite, color.Bold)
	yellow = color.New(color.FgHiYellow, color.Bold)
	red    = color.New(color.FgHiRed, color.Bold)

	namecol = color.New(color.FgHiGreen, color.Bold)
	timecol = color.New(color.FgHiMagenta, color.Bold)
)

/*
const (
	Red        ANSIColor = "\033[1;31m"
	Yellow     ANSIColor = "\033[1;33m"
	Blue       ANSIColor = "\033[1;34m"
	White      ANSIColor = "\033[1;37m"
	ResetColor ANSIColor = "\033[0m"
)
*/
