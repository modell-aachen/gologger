package logtail

import (
	"github.com/pkg/errors"
	"reflect"
	"strings"
	"testing"

	"github.com/modell-aachen/gologger/internal/interfaces"
	"github.com/modell-aachen/gologger/internal/mocks/brokenreporter"
	"github.com/modell-aachen/gologger/internal/mocks/dummydata"
	"github.com/modell-aachen/gologger/internal/mocks/mockreporter"
	"github.com/modell-aachen/gologger/internal/mocks/queue"
)

func TestStoreLogs(t *testing.T) {
	t.Run("Notices when it can not get receiver", func(*testing.T) {
		queueInstance := queue.CreateQueueMock()
		queueInstance.Close()
		mockReporter := &mockreporter.MockReporter{}
		err := TailLogs(queueInstance, mockReporter)
		if !strings.HasPrefix(err.Error(), "Could not connect to queue for logs") {
			t.Error("Did not receive correct error", err.Error())
		}
	})
	t.Run("Closes queue connection on error", func(*testing.T) {
		queueInstance := queue.CreateQueueMock()
		mockReporter := &mockreporter.MockReporter{}
		go queueInstance.MockDelivery(nil, interfaces.LogRow{}, errors.New("MockError"))
		_ = TailLogs(queueInstance, mockReporter)
		if queueInstance.IsClosed() {
			t.Error("Connection is not closed")
		}
	})
	t.Run("Notices when it can not get a delivery", func(*testing.T) {
		queueInstance := queue.CreateQueueMock()
		mockReporter := &mockreporter.MockReporter{}
		go queueInstance.MockDelivery(nil, interfaces.LogRow{}, errors.New("MockError"))
		err := TailLogs(queueInstance, mockReporter)
		if err == nil || !strings.HasPrefix(err.Error(), "Could not receive delivery from queue") {
			t.Error("Got wrong exception", err)
		}
		if errors.Cause(err).Error() != "MockError" {
			t.Error("Got wrong exception origin", err)
		}
	})
	t.Run("Notices when it can not report logs", func(*testing.T) {
		queueInstance := queue.CreateQueueMock()
		mockReporter := &brokenreporter.BrokenReporter{}
		go queueInstance.MockDelivery(nil, interfaces.LogRow{}, nil)
		err := TailLogs(queueInstance, mockReporter)
		if errors.Cause(err).Error() != brokenreporter.Message {
			t.Error("Got wrong exception origin", err)
		}
	})
	t.Run("Reports all logs", func(*testing.T) {
		queueInstance := queue.CreateQueueMock()
		mockReporter := &mockreporter.MockReporter{}
		go TailLogs(queueInstance, mockReporter)
		queueInstance.MockDelivery(nil, dummydata.Row0, nil)
		queueInstance.MockDelivery(nil, dummydata.Row1, nil)
		queueInstance.MockDelivery(nil, dummydata.Row0, errors.New("Test finished"))
		if !reflect.DeepEqual(mockReporter.Logs, []interfaces.LogRow{dummydata.Row0, dummydata.Row1}) {
			t.Error("Did not report correct logs", mockReporter.Logs)
		}
	})
}
