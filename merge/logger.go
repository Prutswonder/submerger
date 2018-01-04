package merge

import (
	"fmt"
)

// Logger is an interface for print methods.
type Logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type loggerImpl struct {
}

// NewLogger instantiates a new Logger that writes to stdout.
func NewLogger() Logger {
	return &loggerImpl{}
}

func (l loggerImpl) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (l loggerImpl) Println(v ...interface{}) {
	fmt.Println(v...)
}
