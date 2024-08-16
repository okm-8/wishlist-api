package config

import (
	"api/internal/model/log"
	"io"
)

type LoggerConfig struct {
	minLevel log.Level
	output   io.Writer
	colored  bool
}

func NewLoggerConfig(minLevel log.Level, output io.Writer, colored bool) *LoggerConfig {
	return &LoggerConfig{
		minLevel: minLevel,
		output:   output,
		colored:  colored,
	}
}

func (config *LoggerConfig) MinLevel() log.Level {
	return config.minLevel
}

func (config *LoggerConfig) Output() io.Writer {
	return config.output
}

func (config *LoggerConfig) Colored() bool {
	return config.colored
}
