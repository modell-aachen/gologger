package logstore

import (
	"log"

	"github.com/robfig/cron"

	"github.com/modell-aachen/gologger/internal/consoleprinter"
	"github.com/modell-aachen/gologger/internal/logstore"
	"github.com/modell-aachen/gologger/internal/postgres"
	"github.com/modell-aachen/gologger/internal/rabbitmq"
)

func Run() {
	queueInstance, err := rabbitmq.CreateInstance()
	if err != nil {
		log.Fatalf("Could not create rabbitmq instance: %+v", err)
	}
	defer queueInstance.Close()

	postgresStore, err := postgres.CreateInstance()
	if err != nil {
		log.Fatalf("Could not connect to postgres: %+v", err)
	}
	defer postgresStore.Close()

	cronRunner := cron.New()
	logstore.ScheduleCleanLogs(cronRunner, postgresStore)

	go func() {
		err := logstore.ReadLogs(queueInstance, postgresStore)
		log.Fatal(err)
	}()

	emergencyLogger := consoleprinter.CreateInstance(true)
	err = logstore.StoreLogs(queueInstance, postgresStore, emergencyLogger)
	if err != nil {
		log.Fatal(err)
	}
}
