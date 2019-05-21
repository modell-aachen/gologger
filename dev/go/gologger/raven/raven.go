package raven

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strings"

	"github.com/modell-aachen/gologger/interfaces"

	"github.com/getsentry/raven-go"
)

type extra interfaces.LogGeneric

func (_ extra) Class() string {
	return "extra"
}

type ravenInstance struct {
}

func GetReporterInstance() (interfaces.LogReporter, error) {
	return ravenInstance{}, nil
}

func (instance ravenInstance) Send(metadata interfaces.LogMetadata, logRow interfaces.LogRow) (err error) {
	if metadata["report"] == "true" {
		message, tags, context, err := createRavenPacket(metadata, logRow)
		if err != nil {
			return errors.Wrapf(err, "Unable to create packet to send for %+v", logRow)
		}
		fmt.Printf("raven %s %+v\n %v", message, context, tags)
		//_ = raven.CaptureMessageAndWait(message, tags, context...)
	}
	return err
}

func createRavenPacket(metadata interfaces.LogMetadata, logRow interfaces.LogRow) (message string, tags map[string]string, context []raven.Interface, err error) {
	if caller, ok := logRow.Misc["caller"]; ok {
		message = caller + ": "
	}
	if len(logRow.Fields) > 0 {
		message = message + "| " + strings.Join(logRow.Fields, " | ") + " |"
	}
	if message == "" {
		message = "(no details available)"
	}
	if len(metadata["tags"]) > 0 {
		err = json.Unmarshal([]byte(metadata["tags"]), &tags)
		delete(metadata, "tags")
	}
	context = make([]raven.Interface, 1)
	context[0] = extra(logRow.Misc)

	return message, tags, context, err
}
