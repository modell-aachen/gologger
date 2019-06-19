package logtail

import (
	"github.com/pkg/errors"

	"github.com/modell-aachen/gologger/internal/interfaces"
)

func TailLogs(queueInstance interfaces.QueueInstance, reporter interfaces.LogReporter) error {
	receiver, err := queueInstance.GetReceiver("")
	if err != nil {
		return errors.Wrap(err, "Could not connect to queue for logs")
	}
	defer receiver.Close()

	for {
		_, logRow, err := receiver.GetDelivery()
		if err != nil {
			return errors.Wrap(err, "Could not receive delivery from queue")
		}

		err = reporter.Send(nil, logRow)
		if err != nil {
			return errors.Wrap(err, "Could not report logs")
		}
	}
}
