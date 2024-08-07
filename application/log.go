package application

import (
	"fmt"
	"io"
	"log"
	"os"
)

var (
	logPrefix   = fmt.Sprintf("[%d] ", os.Getpid())
	logFlags    = log.LstdFlags
	logFileMode = os.FileMode(0644)
)

func resetStdLogger(w io.Writer) {
	log.SetPrefix(logPrefix)
	log.SetFlags(logFlags)
	log.SetOutput(w)
}

// Logger ...
type Logger struct {
	*log.Logger
}

func NewLogger(w io.Writer, reset bool) *Logger {
	l := &Logger{log.New(w, logPrefix, logFlags)}
	if reset {
		resetStdLogger(w)
	}
	return l
}

func NewFileLogger(filename string, reset bool) (*Logger, error) {
	flags := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	w, err := os.OpenFile(filename, flags, logFileMode)
	if err == nil {
		if reset {
			os.Stdout = w
			os.Stderr = w
		}
		return NewLogger(w, reset), nil
	}
	return nil, err
}
