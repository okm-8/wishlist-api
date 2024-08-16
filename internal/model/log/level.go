package log

import "strings"

type Level int

const (
	NotSet Level = iota - 1
	Debug
	Info
	Warning
	Error
)

func ParseLevel(level string) Level {
	level = strings.TrimSpace(level)
	level = strings.ToUpper(level)

	switch level {
	case "DEBUG":
		return Debug
	case "INFO":
		return Info
	case "WARNING", "WARN":
		return Warning
	case "ERROR":
		return Error
	default:
		return NotSet
	}
}

func (level Level) String() string {
	switch level {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warning:
		return "WARNING"
	case Error:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}
