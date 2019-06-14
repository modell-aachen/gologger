package logreport

import (
	"log"

	"github.com/modell-aachen/gologger/internal/logreport"
	"github.com/modell-aachen/gologger/internal/rabbitmq"
	"github.com/modell-aachen/gologger/internal/raven"
)

func Run() {
	queueInstance, err := rabbitmq.CreateInstance()
	if err != nil {
		log.Fatalf("Could not create rabbitmq instance: %+v", err)
	}
	defer queueInstance.Close()

	ravenInstance, err := raven.GetReporterInstance()
	if err != nil {
		log.Fatalf("Could not create raven instance: %+v", err)
	}

	err = logreport.ReportLogs(queueInstance, ravenInstance)
	if err != nil {
		log.Fatal(err)
	}
}
