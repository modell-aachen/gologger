package logstore

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/modell-aachen/gologger/internal/interfaces"
	"github.com/modell-aachen/gologger/internal/mocks/dummydata"
	"github.com/modell-aachen/gologger/internal/mocks/mockreporter"
	"github.com/modell-aachen/gologger/internal/mocks/queue"
	"github.com/modell-aachen/gologger/internal/mocks/store"
)

type jobRunnerDummy struct {
	time     string
	callback func()
	started  bool
}

func (instance *jobRunnerDummy) AddFunc(time string, callback func()) error {
	instance.time = time
	instance.callback = callback
	return nil
}
func (instance *jobRunnerDummy) Start() {
	instance.started = true
}

func TestCleanLogs(t *testing.T) {
	jobRunnerInstance := jobRunnerDummy{}
	storeInstance := store.CreateStoreMock()
	logRow := interfaces.LogRow{}
	storeInstance.Store(logRow)

	t.Run("setup schedule", func(t *testing.T) {
		ScheduleCleanLogs(&jobRunnerInstance, storeInstance)
		if !jobRunnerInstance.started {
			t.Errorf("Cron was not started")
		}
		if jobRunnerInstance.time == "" {
			t.Errorf("No time was added to cron")
		}
		if jobRunnerInstance.callback == nil {
			t.Errorf("Missing callback for cron")
		}
	})

	t.Run("execute callback", func(*testing.T) {
		if len(storeInstance.Stored) != 1 {
			t.Errorf("Set up failed, nothing stored %+v", storeInstance)
		}
		jobRunnerInstance.callback()
		if len(storeInstance.Stored) != 0 {
			t.Errorf("Cleanup callback failed")
		}
	})

	t.Run("notices when callback fails", func(*testing.T) {
		storeInstance.Close()

		defer func() {
			msg := fmt.Sprintf("%s", recover())
			if !strings.HasPrefix(msg, "Could not clean logs") {
				t.Error("Callback paniced with wrong message", msg)
			}
		}()

		jobRunnerInstance.callback()
		t.Error("Callback should have paniced")
	})
}

func TestReadLogsFromStore(t *testing.T) {
	t.Run("replies to requests", func(*testing.T) {
		storeInstance := store.CreateStoreMock()
		storeInstance.Fill()
		queueInstance := queue.CreateQueueMock()

		go ReadLogs(queueInstance, storeInstance)

		queueInstance.MockRpcRequest("c_id", true, time.Unix(0, 0), dummydata.Row1.Time, dummydata.Source, dummydata.Levels)

		published := <-queueInstance.GetRpcChannel("rpc_channel").Publishings
		if published.Key != "rpc_channel" {
			t.Error("Published to wrong channel (mock error?)", published.Key)
		}
		if published.Msg.CorrelationId != "c_id" {
			t.Error("Correlation id is wrong", published.Msg.CorrelationId)
		}

		var received []interfaces.LogRow
		json.Unmarshal(published.Msg.Body, &received)

		expected := []interfaces.LogRow{
			dummydata.Row0,
			dummydata.Row1,
		}

		if !reflect.DeepEqual(received, expected) {
			t.Error("Received wrong logs", received)
		}
	})
	t.Run("Notices when it can not connect to the queue", func(*testing.T) {
		storeInstance := store.CreateStoreMock()
		queueInstance := queue.CreateQueueMock()
		queueInstance.Close()
		err := ReadLogs(queueInstance, storeInstance)
		if !strings.HasPrefix(err.Error(), "Could not connect to queue for rpc") {
			t.Error("Did not receive correct error", err.Error())
		}
	})
	t.Run("Notices when it can not connect to the store", func(*testing.T) {
		storeInstance := store.CreateStoreMock()
		storeInstance.Close()
		queueInstance := queue.CreateQueueMock()
		go queueInstance.MockRpcRequest("c_id", true, time.Unix(0, 0), dummydata.Row2.Time, dummydata.Source, dummydata.Levels)
		err := ReadLogs(queueInstance, storeInstance)
		if !strings.HasPrefix(err.Error(), "Could not read logs") {
			t.Error("Did not receive correct error", err.Error())
		}
	})
	t.Run("Notices when it can not read rpc requests", func(*testing.T) {
		storeInstance := store.CreateStoreMock()
		queueInstance := queue.CreateQueueMock()
		go queueInstance.MockMalformedRpcRequest("c_id", true)
		err := ReadLogs(queueInstance, storeInstance)
		if !strings.HasPrefix(err.Error(), "Could not receive rpc request") {
			t.Error("Did not receive correct error", err.Error())
		}
	})
	t.Run("Notices when it can not reply to rpc requests", func(*testing.T) {
		storeInstance := store.CreateStoreMock()
		queueInstance := queue.CreateQueueMock()
		go func() {
			queueInstance.MockRpcRequest("c_id", false, time.Unix(0, 0), dummydata.Row2.Time, dummydata.Source, dummydata.Levels)

			_ = <-queueInstance.GetRpcChannel("rpc_channel").Publishings
		}()
		err := ReadLogs(queueInstance, storeInstance)
		if !strings.HasPrefix(err.Error(), "Could not reply to read request") {
			t.Error("Did not receive correct error", err.Error())
		}
	})
}

func TestStoreLogs(t *testing.T) {
	mockReporter := &mockreporter.MockReporter{}
	t.Run("Notices when it can not store logs", func(*testing.T) {
		storeInstance := store.CreateStoreMock()
		storeInstance.Close()
		queueInstance := queue.CreateQueueMock()
		go queueInstance.MockDelivery(interfaces.LogMetadata{}, interfaces.LogRow{}, nil)
		err := StoreLogs(queueInstance, storeInstance, mockReporter)
		if !strings.HasPrefix(err.Error(), "Could not store logs in postgres") {
			t.Error("Could not store logs in postgres", err.Error())
		}
	})
	t.Run("Uses fallback when it can not store logs", func(*testing.T) {
		emergencyReporter := &mockreporter.MockReporter{}
		storeInstance := store.CreateStoreMock()
		storeInstance.Close()
		queueInstance := queue.CreateQueueMock()
		go queueInstance.MockDelivery(interfaces.LogMetadata{}, dummydata.Row1, nil)
		StoreLogs(queueInstance, storeInstance, emergencyReporter)
		if !reflect.DeepEqual(emergencyReporter.Logs, []interfaces.LogRow{dummydata.Row1}) {
			t.Error("Fallback was not called correctly", emergencyReporter.Logs)
		}
	})
	t.Run("Notices when it can not read deliveries", func(*testing.T) {
		storeInstance := store.CreateStoreMock()
		queueInstance := queue.CreateQueueMock()
		queueInstance.Close()
		err := StoreLogs(queueInstance, storeInstance, mockReporter)
		if !strings.HasPrefix(err.Error(), "Could not connect to queue for logs") {
			t.Error("Did not receive correct error", err.Error())
		}
	})
	t.Run("Closes the queue connection when finished", func(*testing.T) {
		storeInstance := store.CreateStoreMock()
		queueInstance := queue.CreateQueueMock()
		storeInstance.Close()
		go queueInstance.MockDelivery(nil, interfaces.LogRow{}, errors.New("MockError"))
		err := StoreLogs(queueInstance, storeInstance, mockReporter)
		if !storeInstance.IsClosed() {
			t.Error("Did not close the queue connection", err.Error())
		}
	})
}
