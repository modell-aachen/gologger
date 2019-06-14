package logdump

import (
	"github.com/pkg/errors"
	"time"

	"github.com/modell-aachen/gologger/internal/interfaces"
)

func DumpLogs(storeInstance interfaces.LogStore, reporter interfaces.LogReporter, startTime time.Time, endTime time.Time, source interfaces.SourceString, levels []interfaces.LevelString) error {
	logRows, err := storeInstance.Read(startTime, endTime, source, levels)
	if err != nil {
		return errors.Wrap(err, "Could not read logs")
	}

	for _, logRow := range logRows {
		err = reporter.Send(nil, logRow)
		if err != nil {
			return errors.Wrap(err, "Could not print logs")
		}
	}

	return nil
}
