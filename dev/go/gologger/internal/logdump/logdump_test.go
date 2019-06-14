package logdump

import (
	"github.com/pkg/errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/modell-aachen/gologger/internal/interfaces"
	"github.com/modell-aachen/gologger/internal/mocks/brokenreporter"
	"github.com/modell-aachen/gologger/internal/mocks/dummydata"
	"github.com/modell-aachen/gologger/internal/mocks/mockreporter"
	"github.com/modell-aachen/gologger/internal/mocks/store"
)

func TestStoreLogs(t *testing.T) {
	t.Run("Notices when it can not read logs", func(*testing.T) {
		mockReporter := &mockreporter.MockReporter{}
		storeInstance := store.CreateStoreMock()
		storeInstance.Close()
		err := DumpLogs(storeInstance, mockReporter, time.Unix(0, 0), time.Unix(0, 0), dummydata.Source, dummydata.Levels)
		if !strings.HasPrefix(err.Error(), "Could not read logs") {
			t.Error("Got wrong exception", err.Error())
		}
	})
	t.Run("Notices when it can not report logs", func(*testing.T) {
		mockReporter := &brokenreporter.BrokenReporter{}
		storeInstance := store.CreateStoreMock()
		storeInstance.Fill()
		err := DumpLogs(storeInstance, mockReporter, time.Unix(0, 0), dummydata.Row2.Time, dummydata.Source, dummydata.Levels)
		if err == nil || errors.Cause(err).Error() != brokenreporter.Message {
			t.Error("Got wrong exception", err)
		}
	})
	t.Run("Reports all logs", func(*testing.T) {
		mockReporter := &mockreporter.MockReporter{}
		storeInstance := store.CreateStoreMock()
		storeInstance.Fill()
		_ = DumpLogs(storeInstance, mockReporter, time.Unix(0, 0), dummydata.Row2.Time, dummydata.Source, dummydata.Levels)
		if !reflect.DeepEqual(mockReporter.Logs, []interfaces.LogRow{dummydata.Row0, dummydata.Row1, dummydata.Row2}) {
			t.Error("Did not report correct logs", mockReporter.Logs)
		}
	})
}
