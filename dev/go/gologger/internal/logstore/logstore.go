package logstore

import (
	"github.com/pkg/errors"

	"github.com/robfig/cron"

	"github.com/modell-aachen/gologger/internal/interfaces"
)

type jobRunner interface {
	AddFunc(string, func()) error
	Start()
}

func ScheduleCleanLogs(jobRunnerInstance jobRunner, store interfaces.LogStore) {
	jobRunnerInstance.AddFunc(
		"5 5 5 * * *",
		CreateFuncJob(store),
	)
	jobRunnerInstance.Start()
}

func CreateFuncJob(store interfaces.LogStore) func() {
	return (cron.FuncJob)(func() {
		err := store.CleanUp()
		if err != nil {
			panic(errors.Wrapf(err, "Could not clean logs"))
		}
	})
}

func ReadLogs(queueInstance interfaces.QueueInstance, store interfaces.LogStore) error {
	rpc, err := queueInstance.GetRpcReceiver()
	if err != nil {
		return errors.Wrapf(err, "Could not connect to queue for rpc")
	}
	defer rpc.Close()

	for {
		request, startTime, endTime, source, levels, err := rpc.GetRequest()
		if err != nil {
			return errors.Wrapf(err, "Could not receive rpc request")
		}

		fields, err := store.Read(startTime, endTime, source, levels)
		if err != nil {
			return errors.Wrapf(err, "Could not read logs")
		}
		err = request.Reply(fields)
		if err != nil {
			return errors.Wrapf(err, "Could not reply to read request")
		}
	}
}

func StoreLogs(queueInstance interfaces.QueueInstance, store interfaces.LogStore, emergencyLogger interfaces.LogReporter) error {
	receiver, err := queueInstance.GetReceiver("log_store_go")
	if err != nil {
		return errors.Wrapf(err, "Could not connect to queue for logs")
	}
	defer receiver.Close()

	for {
		_, logRow, err := receiver.GetDelivery()
		if err != nil {
			return errors.Wrapf(err, "Could not receive delivery from rabbitmq")
		}

		err = store.Store(logRow)
		if err != nil {
			emergencyLogger.Send(nil, logRow)
			return errors.Wrapf(err, "Could not store logs in postgres")
		}
	}
}
