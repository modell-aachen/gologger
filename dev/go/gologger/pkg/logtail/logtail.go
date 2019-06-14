package logtail

import (
	"log"
	"os"

	"github.com/modell-aachen/gologger/internal/consoleprinter"
	internallogtail "github.com/modell-aachen/gologger/internal/logtail"
	"github.com/modell-aachen/gologger/internal/rabbitmq"
)

func Run() {
	fi, _ := os.Stdout.Stat()
	piped := (fi.Mode() & os.ModeCharDevice) == 0

	queueInstance, err := rabbitmq.CreateInstance()
	if err != nil {
		log.Fatalf("Could not create rabbitmq instance: %+v", err)
	}
	defer queueInstance.Close()

	consolePrinter := consoleprinter.CreateInstance(piped)

	err = internallogtail.TailLogs(queueInstance, consolePrinter)
	if err != nil {
		log.Fatal(err)
	}
}
