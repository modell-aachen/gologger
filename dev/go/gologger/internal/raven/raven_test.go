package raven

import (
	"reflect"
	"testing"

	"github.com/modell-aachen/gologger/internal/interfaces"
	"github.com/modell-aachen/gologger/internal/mocks/dummydata"

	"github.com/getsentry/raven-go"
)

func TestStoreLogs(t *testing.T) {
	raven.SetDSN("") // prevent it from sending anything
	t.Run("Creates the raven packet correctly", func(*testing.T) {
		t.Run("with data", func(*testing.T) {
			metadata := interfaces.LogMetadata{
				"tags": "{\"version\":\"test\"}",
			}
			message, tags, context, err := createRavenPacket(metadata, dummydata.Row1)
			if err != nil {
				t.Error(err)
			}
			if message != "Mr. R. A.: | field a | field b |" {
				t.Error("Incorrect message", message)
			}
			if !reflect.DeepEqual(map[string]string{"version": "test"}, tags) {
				t.Error("Wrong tags", tags)
			}
			expectedContext := []raven.Interface{
				extra(dummydata.Row1.Misc),
			}
			if !reflect.DeepEqual(expectedContext, context) {
				t.Error("Wrong context", context)
			}
		})
		t.Run("without a caller", func(*testing.T) {
			metadata := interfaces.LogMetadata{}
			message, _, _, _ := createRavenPacket(metadata, dummydata.Row0)
			if message != "| early | entry |" {
				t.Error("Incorrect message", message)
			}
		})
		t.Run("with no data", func(*testing.T) {
			metadata := interfaces.LogMetadata{}
			message, tags, context, err := createRavenPacket(metadata, interfaces.LogRow{})
			if err != nil {
				t.Error(err)
			}
			if message != "(no details available)" {
				t.Error("Incorrect message", message)
			}
			if len(tags) != 0 {
				t.Error("Wrong tags", tags)
			}
			expectedContext := []raven.Interface{
				extra(nil),
			}
			if !reflect.DeepEqual(expectedContext, context) {
				t.Error("Wrong context", context)
			}
		})
	})
	t.Run("Creates a raven instance", func(*testing.T) {
		raven, _ := GetReporterInstance()
		_, ok := raven.(ravenInstance)
		if !ok {
			t.Fail()
		}
	})
	t.Run("Sends when metadata asks it to", func(*testing.T) {
		instance, _ := GetReporterInstance()
		err := instance.Send(interfaces.LogMetadata{"report": "true"}, dummydata.Row1)
		if err != nil {
			t.Fail()
		}
	})
	t.Run("Reports errors when attempting to send", func(*testing.T) {
		instance, _ := GetReporterInstance()
		err := instance.Send(interfaces.LogMetadata{"report": "true", "tags": "not json"}, dummydata.Row1)
		if err == nil {
			t.Fail()
		}
	})
	t.Run("Type 'extra' has class 'extra' for raven packet", func(*testing.T) {
		if extra(map[string]string{"key": "value"}).Class() != "extra" {
			t.Fail()
		}
	})
}
