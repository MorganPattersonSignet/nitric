package noninteractive

import (
	"io"
	"strings"

	"github.com/charmbracelet/log"
)

type NonInteractiveOutput struct {
	logger *log.Logger
}

func (n *NonInteractiveOutput) Write(bytes []byte) (int, error) {
	msg := string(bytes)
	msg, _ = strings.CutSuffix(msg, "\n")
	n.logger.Print(msg)
	return len(bytes), nil
}

func NewNonInterativeOutput(output io.Writer) io.Writer {
	logger := log.New(output)
	logger.SetPrefix("")

	return &NonInteractiveOutput{
		logger: logger,
	}
}
