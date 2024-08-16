package output

import (
	"errors"
	"io"
	"os"
	"strings"
)

const (
	Stdout = "stdout"
	Stderr = "stderr"
)

var ErrUnsupportedOutput = errors.New("unsupported output")

func New(name string) (io.Writer, error) {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)

	switch name {
	case Stdout:
		return os.Stdout, nil
	case Stderr:
		return os.Stderr, nil
	}

	if strings.HasPrefix(name, "file://") {
		return os.Create(strings.TrimPrefix(name, "file://"))
	}

	return nil, ErrUnsupportedOutput
}
