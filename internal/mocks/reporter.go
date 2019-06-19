package mocks

import (
	"errors"

	"github.com/modell-aachen/gologger/internal/interfaces"
	"github.com/modell-aachen/gologger/internal/mocks/queue"
)

type MockReporter struct {
	Reported []queue.MockDelivery
	hasError bool
}

func (instance *MockReporter) Send(metadata interfaces.LogMetadata, logRow interfaces.LogRow) error {
	if instance.hasError {
		return errors.New("Error mocked")
	}
	instance.Reported = append(instance.Reported, queue.MockDelivery{metadata, logRow, nil})
	return nil
}

func (instance *MockReporter) EnterErrorState() {
	instance.hasError = true
}

var _ interfaces.LogReporter = (*MockReporter)(nil)
