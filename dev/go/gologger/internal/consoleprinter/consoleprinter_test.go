package consoleprinter

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/modell-aachen/gologger/internal/interfaces"
	"github.com/modell-aachen/gologger/internal/mocks/dummydata"
)

func TestStoreLogs(t *testing.T) {
	t.Run("Creates piped instances", func(*testing.T) {
		instance := CreateInstance(true)
		if instance.rules != &formatRulesPlain {
			t.Error("Did not create plain rules", *instance)
		}
	})
	t.Run("Creates console instances", func(*testing.T) {
		instance := CreateInstance(false)
		if instance.rules != &formatRulesFormatted {
			t.Error("Did not create formatted rules", *instance)
		}

	})
	t.Run("Formats Row1 correctly", func(*testing.T) {
		instance := CreateInstance(true)
		output := instance.formatForOutput(dummydata.Row1)

		expected := []interface{}{
			instance.rules.start,
			fmt.Sprint(dummydata.Row1.Time),
			instance.rules.separator,
			fmt.Sprint(dummydata.Level),
			instance.rules.separator,
			fmt.Sprint(dummydata.Source),
			instance.rules.separator,
			"caller: Mr. R. A.",
			instance.rules.separator,
			"key a: value a",
			instance.rules.separator,
			"key b: value b",
			instance.rules.separator,
			dummydata.Row1.Fields[0],
			instance.rules.separator,
			dummydata.Row1.Fields[1],
			instance.rules.end,
		}
		if !reflect.DeepEqual(output, expected) {
			t.Error("Did not receive expected result", output, expected)
		}
	})
	t.Run("Sorts log keys correctly", func(*testing.T) {
		var logs = interfaces.LogGeneric{
			"X": "Last",
			"B": "Second",
			"A": "First",
		}
		output := getSortedLogKeys(logs)
		expected := []string{"A", "B", "X"}
		if !reflect.DeepEqual(output, expected) {
			t.Error("Did not sort as expected", output, expected)
		}
	})
}
