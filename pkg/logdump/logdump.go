package logdump

import (
	"log"
	"os"
	"time"

	"github.com/modell-aachen/gologger/internal/consoleprinter"
	"github.com/modell-aachen/gologger/internal/interfaces"
	internallogdump "github.com/modell-aachen/gologger/internal/logdump"
	"github.com/modell-aachen/gologger/internal/postgres"
)

func Run() {
	fi, _ := os.Stdout.Stat()
	piped := (fi.Mode() & os.ModeCharDevice) == 0

	dbInstance, err := postgres.CreateInstance()
	if err != nil {
		log.Fatalf("Could not create postgres instance: %+v", err)
	}
	defer dbInstance.Close()

	consolePrinter := consoleprinter.CreateInstance(piped)

	startTime := time.Unix(0, 0)
	endTime := time.Now()
	source := interfaces.SourceString("")
	levels := []interfaces.LevelString{
		interfaces.LevelString("notice"),
		interfaces.LevelString("event"),
		interfaces.LevelString("debug"),
		interfaces.LevelString("info"),
		interfaces.LevelString("warning"),
		interfaces.LevelString("error"),
		interfaces.LevelString("fatal"),
	}
	err = internallogdump.DumpLogs(dbInstance, consolePrinter, startTime, endTime, source, levels)

	if err != nil {
		log.Fatal("Could not dump logs", err)
	}
}
