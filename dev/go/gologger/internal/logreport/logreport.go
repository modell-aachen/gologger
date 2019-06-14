package logreport

import (
	"github.com/pkg/errors"

	"github.com/modell-aachen/gologger/internal/interfaces"
)

func ReportLogs(queueInstance interfaces.QueueInstance, reporterInstance interfaces.LogReporter) error {
	receiver, err := queueInstance.GetReceiver("")
	if err != nil {
		return errors.Wrapf(err, "Could not connect to queue for logs")
	}
	defer receiver.Close()

	for {
		LogMetadata, logRow, err := receiver.GetDelivery()
		if err != nil {
			return errors.Wrapf(err, "Could not receive delivery from queue")
		}

		err = reporterInstance.Send(LogMetadata, logRow)
		if err != nil {
			return errors.Wrapf(err, "Could not report delivery")
		}
	}
}
