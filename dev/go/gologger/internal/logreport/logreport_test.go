package logreport

import (
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/modell-aachen/gologger/internal/interfaces"
	"github.com/modell-aachen/gologger/internal/mocks"
	"github.com/modell-aachen/gologger/internal/mocks/queue"
)

func TestReportLogs(t *testing.T) {
	t.Run("Closes connection on error", func(t *testing.T) {
		mockReporter := mocks.CreateReporterMock()
		mockQueue := queue.CreateQueueMock()
		go mockQueue.MockDelivery(nil, interfaces.LogRow{}, errors.New("end of test"))
		err := ReportLogs(mockQueue, mockReporter)
		if errors.Cause(err).Error() != "end of test" {
			t.Errorf("Did not process until end of test: %+v", err)
		}
		receivers := mockQueue.GetReceivers()
		if len(receivers) != 1 {
			t.Errorf("expected exactly one receiver, but got %d", len(receivers))
		}
		if !receivers[0].IsClosed() {
			t.Error("Receiver was not closed")
		}
	})
	t.Run("Pushes received messages to the reporter", func(t *testing.T) {
		mockReporter := mocks.CreateReporterMock()
		mockQueue := queue.CreateQueueMock()

		logMetadata := make(interfaces.LogMetadata)
		logMetadata["test"] = "some value"

		misc := make(interfaces.LogGeneric)
		misc["generic test"] = "generic value"
		logFields := interfaces.LogFields{"a1", "b 1"}
		logRow := interfaces.LogRow{
			time.Unix(554385600, 0),
			"wiki",
			interfaces.LevelString("info"),
			misc,
			logFields,
		}

		go func() {
			mockQueue.MockDelivery(logMetadata, logRow, nil)
			mockQueue.MockDelivery(nil, interfaces.LogRow{}, errors.New("end of test"))
		}()

		err := ReportLogs(mockQueue, mockReporter)
		if errors.Cause(err).Error() != "end of test" {
			t.Errorf("Did not process until end of test: %+v", err)
		}

		if len(mockReporter.Reported) != 1 {
			t.Errorf("Did not report correct number of logs: %d", len(mockReporter.Reported))
		}

		if !reflect.DeepEqual(logRow, mockReporter.Reported[0].LogRow) {
			t.Errorf("Did not report correct log: \n%+v\n vs. \n%+v", logRow, mockReporter.Reported[0].LogRow)
		}
		if !reflect.DeepEqual(logMetadata, mockReporter.Reported[0].Metadata) {
			t.Errorf("Did not report correct metadata: \n%+v\n vs. \n%+v", logMetadata, mockReporter.Reported[0].Metadata)
		}

	})
	t.Run("Reports connection errors", func(t *testing.T) {
		t.Run("with the queue", func(t *testing.T) {
			mockReporter := mocks.CreateReporterMock()
			mockQueue := queue.CreateQueueMock()
			mockQueue.Close()

			err := ReportLogs(mockQueue, mockReporter)
			if errors.Cause(err).Error() != "Mock connection closed" {
				t.Errorf("Expected other error: %+v", err)
			}
		})
		t.Run("with the reporter", func(t *testing.T) {
			mockReporter := mocks.CreateReporterMock()
			mockReporter.EnterErrorState()
			mockQueue := queue.CreateQueueMock()

			go func() {
				mockQueue.MockDelivery(nil, interfaces.LogRow{}, nil)
				mockQueue.MockDelivery(nil, interfaces.LogRow{}, errors.New("end of test"))
			}()

			err := ReportLogs(mockQueue, mockReporter)
			if errors.Cause(err).Error() != "Error mocked" {
				t.Errorf("Expected other error: %+v", err)
			}
		})
	})
}
