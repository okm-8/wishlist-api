package log

import (
	"time"
)

type Label struct {
	key   string
	value any
}

func NewLabel(key string, value any) *Label {
	return &Label{
		key:   key,
		value: value,
	}
}

func (label Label) Key() string {
	return label.key
}

func (label Label) Value() any {
	return label.value
}

type Record struct {
	timestamp time.Time
	level     Level
	message   string
	labels    []*Label
}

func NewRecord(level Level, message string, labels []*Label) Record {
	return Record{
		timestamp: time.Now(),
		level:     level,
		message:   message,
		labels:    labels,
	}
}

func (r Record) Timestamp() time.Time {
	return r.timestamp
}

func (r Record) Level() Level {
	return r.level
}

func (r Record) Message() string {
	return r.message
}

func (r Record) Labels() []*Label {
	return r.labels
}
