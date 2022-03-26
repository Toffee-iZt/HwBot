package logger

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

// New creates new logger instance.
func New(w *Writer, name string) *Logger {
	return &Logger{
		out:  w,
		name: name,
	}
}

// Logger struct.
type Logger struct {
	out  *Writer
	name string
}

// Child ...
func (l *Logger) Child(name string) *Logger {
	return New(l.out, l.name+"::"+name)
}

// Copy ...
func (l *Logger) Copy(name string) *Logger {
	return New(l.out, name)
}

// SetWriter ...
func (l *Logger) SetWriter(w *Writer) {
	l.out = w
}

func (l *Logger) log(pref string, pcol *color.Color, f string, v ...interface{}) {
	if f == "" {
		return
	}

	w := l.out
	w.m.Lock()

	buf := l.printTime(w.b[:0], w.c)
	buf = l.append(buf, pref, pcol)
	if l.name != "" {
		buf = l.append(buf, l.name, namecol)
	}

	buf = append(buf, fmt.Sprintf(f, v...)...)
	buf = append(buf, '\n')

	w.w.Write(buf)
	w.b = buf

	w.m.Unlock()
}

func (l *Logger) printTime(buf []byte, col bool) []byte {
	t := time.Now()
	if past := t.Sub(l.out.l); past > time.Minute {
		l.out.l = t
		return l.append(buf, t.Format("02.01.2006 15:04\n"), timecol)
	}
	return buf
}

func (l *Logger) append(buf []byte, text string, col *color.Color) []byte {
	if l.out.c {
		buf = append(buf, col.Sprint(text)...)
	} else {
		buf = append(buf, text...)
	}
	return append(buf, ' ')
}

// Info ...
func (l *Logger) Info(f string, v ...interface{}) {
	l.log("INFO ", white, f, v...)
}

// Warn ...
func (l *Logger) Warn(f string, v ...interface{}) {
	l.log("WARN ", yellow, f, v...)
}

// Error ...
func (l *Logger) Error(f string, v ...interface{}) {
	l.log("ERROR", red, f, v...)
}
