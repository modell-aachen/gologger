package main

import (
	"log"

	"github.com/modell-aachen/gologger/consoleprinter"
	"github.com/modell-aachen/gologger/interfaces"
	"github.com/modell-aachen/gologger/rabbitmq"
)

func main() {
	queueInstance, err := rabbitmq.CreateInstance()
	if err != nil {
		log.Fatalf("Could not create rabbitmq instance: %+v", err)
	}
	defer queueInstance.Close()

	consolePrinter := consoleprinter.CreateInstance()

	tailLogs(queueInstance, consolePrinter)
}

func tailLogs(queueInstance interfaces.QueueInstance, reporter interfaces.LogReporter) {
	receiver, err := queueInstance.GetReceiver("")
	if err != nil {
		log.Fatalf("Could not connect to queue for logs: %+v", err)
	}
	defer receiver.Close()

	for {
		_, logRow, err := receiver.GetDelivery()
		if err != nil {
			log.Fatalf("Could not receive delivery from rabbitmq: %+v", err)
		}

		reporter.Send(nil, logRow)
	}
}
