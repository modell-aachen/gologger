package main

import (
	"log"
	"time"

	"github.com/modell-aachen/gologger/consoleprinter"
	"github.com/modell-aachen/gologger/interfaces"
	"github.com/modell-aachen/gologger/postgres"
)

func main() {
	dbInstance, err := postgres.CreateInstance()
	if err != nil {
		log.Fatalf("Could not create postgres instance: %+v", err)
	}
	defer dbInstance.Close()

	consolePrinter := consoleprinter.CreateInstance()

	startTime := time.Unix(0, 0)
	endTime := time.Now()
	source := interfaces.SourceString("")
	levels := []interfaces.LevelString{
		interfaces.LevelString("warning"),
	}
	dumpLogs(dbInstance, consolePrinter, startTime, endTime, source, levels)
}

func dumpLogs(storeInstance interfaces.LogStore, reporter interfaces.LogReporter, startTime time.Time, endTime time.Time, source interfaces.SourceString, levels []interfaces.LevelString) {
	logRows, err := storeInstance.Read(startTime, endTime, source, levels)
	if err != nil {
		log.Fatalf("Could not read logs: %+v", err)
	}

	for _, logRow := range logRows {
		reporter.Send(nil, logRow)
	}
}
