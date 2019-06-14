package mockreporter

import (
	"github.com/modell-aachen/gologger/internal/interfaces"
)

type MockReporter struct {
	Logs     []interfaces.LogRow
	Metadata []interfaces.LogMetadata
}

func (instance *MockReporter) Send(metadata interfaces.LogMetadata, log interfaces.LogRow) error {
	instance.Logs = append(instance.Logs, log)
	instance.Metadata = append(instance.Metadata, metadata)

	return nil
}
