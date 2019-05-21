package main

import (
	"log"

	"github.com/modell-aachen/gologger/interfaces"
	"github.com/modell-aachen/gologger/rabbitmq"
	"github.com/modell-aachen/gologger/raven"
)

func main() {
	queueInstance, err := rabbitmq.CreateInstance()
	if err != nil {
		log.Fatalf("Could not create rabbitmq instance: %+v", err)
	}
	defer queueInstance.Close()

	ravenInstance, err := raven.GetReporterInstance()
	if err != nil {
		log.Fatalf("Could not create raven instance: %+v", err)
	}

	reportLogs(queueInstance, ravenInstance)
}

func reportLogs(queueInstance interfaces.QueueInstance, reporterInstance interfaces.LogReporter) {
	receiver, err := queueInstance.GetReceiver("")
	if err != nil {
		log.Fatalf("Could not connect to queue for logs: %+v", err)
	}
	defer receiver.Close()

	for {
		LogMetadata, logRow, err := receiver.GetDelivery()
		if err != nil {
			log.Fatalf("Could not receive delivery from queue: %+v", err)
		}

		err = reporterInstance.Send(LogMetadata, logRow)
		if err != nil {
			log.Fatalf("Could not report delivery: %+v", err)
		}
	}
}
