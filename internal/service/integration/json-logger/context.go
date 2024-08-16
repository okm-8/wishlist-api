package jsonLogger

import "io"

type Context interface {
	LogOutput() io.Writer
	ColoredLog() bool
}
