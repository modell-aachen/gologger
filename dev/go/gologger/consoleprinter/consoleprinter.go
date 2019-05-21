package consoleprinter

import (
	"fmt"
	"sort"

	. "github.com/logrusorgru/aurora"
	"github.com/modell-aachen/gologger/interfaces"
)

type consolePrinter struct{}

func CreateInstance() (instance interfaces.LogReporter) {
	return consolePrinter{}
}

func (instance consolePrinter) Send(metadata interfaces.LogMetadata, logRow interfaces.LogRow) (err error) {
	columns := formatForOutput(logRow)
	fmt.Println(columns...)
	return nil
}

func formatForOutput(logRow interfaces.LogRow) (columns []interface{}) {
	columns = make([]interface{}, 1+2*(3+len(logRow.Misc)+len(logRow.Fields)))
	separator := Bold(" | ")

	columns[0] = Bold("| ")
	columns[1] = logRow.Time.String()
	columns[2] = separator
	columns[3] = string(logRow.Level)
	columns[4] = separator
	columns[5] = string(logRow.Source)

	offset := 6
	for _, key := range getSortedLogKeys(logRow.Misc) {
		columns[offset] = separator
		offset++
		columns[offset] = fmt.Sprintf("%s: %s", key, logRow.Misc[key])
		offset++
	}

	for _, item := range logRow.Fields {
		columns[offset] = separator
		offset++
		columns[offset] = item
		offset++
	}
	columns[offset] = Bold(" |")

	return columns
}

func getSortedLogKeys(logFields interfaces.LogGeneric) (keys []string) {
	keys = make([]string, len(logFields))
	idx := 0
	for key := range logFields {
		keys[idx] = key
		idx += 1
	}
	sort.Strings(keys)
	return keys
}
