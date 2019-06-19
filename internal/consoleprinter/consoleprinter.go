package consoleprinter

import (
	"fmt"
	"sort"

	. "github.com/logrusorgru/aurora"
	"github.com/modell-aachen/gologger/internal/interfaces"
)

type consolePrinter struct {
	rules *formatRules
}

var _ interfaces.LogReporter = (*consolePrinter)(nil)

type formatRules struct {
	start     string
	separator string
	end       string
	levels    map[interfaces.LevelString]string
}

var (
	formatRulesPlain = formatRules{
		start:     "| ",
		separator: " | ",
		end:       " |",
		levels:    map[interfaces.LevelString]string{},
	}
	formatRulesFormatted = formatRules{
		start:     fmt.Sprint(Bold(formatRulesPlain.start)),
		separator: fmt.Sprint(Bold(formatRulesPlain.separator)),
		end:       fmt.Sprint(Bold(formatRulesPlain.end)),
		levels: map[interfaces.LevelString]string{
			interfaces.LevelString("notice"):  fmt.Sprint(BrightGreen("notice")),
			interfaces.LevelString("event"):   fmt.Sprint(BrightGreen("event")),
			interfaces.LevelString("info"):    fmt.Sprint(BrightGreen("info")),
			interfaces.LevelString("debug"):   fmt.Sprint(BrightCyan("debug")),
			interfaces.LevelString("warning"): fmt.Sprint(BrightYellow("warning")),
			interfaces.LevelString("error"):   fmt.Sprint(BrightMagenta("error")),
			interfaces.LevelString("fatal"):   fmt.Sprint(BrightRed("fatal")),
		},
	}
)

func CreateInstance(plain bool) (instance *consolePrinter) {
	if plain {
		instance = &consolePrinter{
			rules: &formatRulesPlain,
		}
	} else {
		instance = &consolePrinter{
			rules: &formatRulesFormatted,
		}
	}
	return instance
}

func (instance *consolePrinter) Send(metadata interfaces.LogMetadata, logRow interfaces.LogRow) (err error) {
	columns := instance.formatForOutput(logRow)
	fmt.Println(columns...)
	return nil
}

func (instance *consolePrinter) formatForOutput(logRow interfaces.LogRow) (columns []interface{}) {
	level, ok := instance.rules.levels[logRow.Level]
	if !ok {
		level = (string)(logRow.Level)
	}
	columns = make([]interface{}, 1+2*(3+len(logRow.Misc)+len(logRow.Fields)))

	columns[0] = instance.rules.start
	columns[1] = logRow.Time.String()
	columns[2] = instance.rules.separator
	columns[3] = level
	columns[4] = instance.rules.separator
	columns[5] = string(logRow.Source)

	offset := 6
	for _, key := range getSortedLogKeys(logRow.Misc) {
		columns[offset] = instance.rules.separator
		offset++
		columns[offset] = fmt.Sprintf("%s: %s", key, logRow.Misc[key])
		offset++
	}

	for _, item := range logRow.Fields {
		columns[offset] = instance.rules.separator
		offset++
		columns[offset] = item
		offset++
	}
	columns[offset] = instance.rules.end

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
