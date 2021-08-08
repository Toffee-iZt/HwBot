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
	last time.Time
	out  *Writer
	name string
}

// Child ...
func (l *Logger) Child(name string) *Logger {
	newl := *l
	newl.name = newl.name + "::" + name
	return &newl
}

// Copy ...
func (l *Logger) Copy(name string) *Logger {
	newl := *l
	newl.name = name
	return &newl
}

// SetWriter ...
func (l *Logger) SetWriter(w *Writer) {
	l.out = w
}

// Writer ...
func (l *Logger) Writer() *Writer {
	return l.out
}

func (l *Logger) log(pref string, pcol *color.Color, f string, v ...interface{}) {
	if f == "" {
		return
	}

	w := l.out
	w.m.Lock()

	buf := l.printTime(w.b[:0], w.c)

	if w.c {
		buf = append(buf, pcol.Sprint(pref)...)
	} else {
		buf = append(buf, pref...)
	}
	buf = append(buf, ' ')

	if l.name != "" {
		if w.c {
			buf = append(buf, namecol.Sprint(l.name)...)
		} else {
			buf = append(buf, l.name...)
		}
		buf = append(buf, ' ')
	}

	buf = append(buf, fmt.Sprintf(f, v...)...)
	buf = append(buf, '\n')

	w.w.Write(buf)
	w.b = buf

	w.m.Unlock()
}

func (l *Logger) printTime(buf []byte, col bool) []byte {
	t := time.Now()
	_, nm, _ := t.Clock()
	if _, m, _ := l.last.Clock(); m != nm || t.Sub(l.last) > time.Minute {
		l.last = t
		tstr := t.Format("02.01.2006 15:04\n")
		if col {
			tstr = timecol.Sprint(tstr)
		}
		return append(buf, tstr...)
	}
	return buf
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
