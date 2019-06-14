package brokenreporter

import (
	"errors"
	"github.com/modell-aachen/gologger/internal/interfaces"
)

var Message string = "Can not report logs"

type BrokenReporter struct{}

func (instance *BrokenReporter) Send(metadata interfaces.LogMetadata, logRow interfaces.LogRow) error {
	return errors.New(Message)
}
