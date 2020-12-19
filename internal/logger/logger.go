package logger

import (
	"fmt"
	"io"
	"os"
)

type logger struct {
}

func New() *logger {
	return &logger{}
}

func (l logger) Info(format string, a ...interface{}) {
	_, _ = io.WriteString(os.Stdout, fmt.Sprintf(format, a...))
	_, _ = io.WriteString(os.Stdout, "\n")
}

func (l logger) Error(format string, a ...interface{}) {
	_, _ = io.WriteString(os.Stderr, fmt.Sprintf(format, a...))
	_, _ = io.WriteString(os.Stdout, "\n")
}

func (l logger) Fatal(format string, a ...interface{}) {
	l.Error(format, a...)
	os.Exit(1)
}
