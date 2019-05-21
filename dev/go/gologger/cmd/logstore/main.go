package main

import (
	"fmt"
	"log"

	"github.com/robfig/cron"

	"github.com/modell-aachen/gologger/interfaces"
	"github.com/modell-aachen/gologger/postgres"
	"github.com/modell-aachen/gologger/rabbitmq"
)

type jobRunner interface {
	AddFunc(string, func()) error
	Start()
}

func main() {
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
	scheduleCleanLogs(cronRunner, postgresStore)

	go readLogs(queueInstance, postgresStore)
	storeLogs(queueInstance, postgresStore)
}

func scheduleCleanLogs(jobRunnerInstance jobRunner, store interfaces.LogStore) {
	jobRunnerInstance.AddFunc(
		"5 5 5 * * *",
		func() {
			err := store.CleanUp()
			if err != nil {
				log.Printf("Could not clean logs: %+v", err)
			}
		},
	)
	jobRunnerInstance.Start()
}

func readLogs(queueInstance interfaces.QueueInstance, store interfaces.LogStore) {
	rpc, err := queueInstance.GetRpcReceiver()
	if err != nil {
		log.Fatalf("Could not connect to queue for rpc: %+v", err)
	}
	defer rpc.Close()

	for {
		delivery, startTime, endTime, source, levels, err := rpc.GetRequest()
		if err != nil {
			log.Fatalf("Could not receive rpc request: %+v", err)
		}

		fields, err := store.Read(startTime, endTime, source, levels)
		if err != nil {
			log.Fatalf("Could not read logs: %+v", err)
		}
		err = rpc.ReplyToRequest(delivery, fields)
		if err != nil {
			log.Fatalf("Could not reply to read request: %+v", err)
		}
	}
}

func storeLogs(queueInstance interfaces.QueueInstance, store interfaces.LogStore) {
	receiver, err := queueInstance.GetReceiver("log_store_go")
	if err != nil {
		log.Fatalf("Could not connect to queue for logs: %+v", err)
	}
	defer receiver.Close()

	for {
		_, logRow, err := receiver.GetDelivery()
		if err != nil {
			log.Fatalf("Could not receive delivery from rabbitmq: %+v", err)
		}

		err = store.Store(logRow)
		if err != nil {
			fmt.Printf("%s %s (%s): %v | %v\n", logRow.Time, logRow.Level, logRow.Source, logRow.Misc, logRow.Fields)
			log.Fatalf("Could not store logs in postgres: %+v", err)
		}
	}
}
