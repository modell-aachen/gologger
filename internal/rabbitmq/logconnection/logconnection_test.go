package logconnection

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/modell-aachen/gologger/internal/interfaces"
	"github.com/modell-aachen/gologger/internal/mocks/deliverysupplier"
	"github.com/modell-aachen/gologger/internal/mocks/dummydata"
	"github.com/modell-aachen/gologger/internal/mocks/mockdelivery"
)

type mockChannel struct {
	isClosed bool
}

func (instance *mockChannel) Close() error {
	instance.isClosed = true
	return nil
}

func createTestInstance() (*LogConnection, *mockChannel, deliverysupplier.DeliverySupplier) {
	channel := &mockChannel{}
	supplier := deliverysupplier.CreateDeliverySupplier()
	connection := CreateLogConnection(
		channel,
		supplier,
	)

	return connection, channel, supplier
}

func TestGetDelivery(t *testing.T) {
	t.Run("Closes acks the delivery", func(t *testing.T) {
		instance, _, supplier := createTestInstance()
		delivery := mockdelivery.CreateMockDelivery("test", "123", nil, true)
		go func() {
			supplier <- delivery
		}()
		instance.GetDelivery()
		if !delivery.IsAcked() {
			t.Fail()
		}
	})
	t.Run("Notices when it can not ack the delivery", func(t *testing.T) {
		instance, _, supplier := createTestInstance()
		delivery := mockdelivery.CreateMockDelivery("test", "123", nil, false)
		go func() {
			supplier <- delivery
		}()
		_, _, err := instance.GetDelivery()
		if !strings.HasPrefix(err.Error(), "Could not ACK delivery") {
			t.Fail()
		}
	})
	t.Run("Notices when it can not unmarshal the message", func(t *testing.T) {
		instance, _, supplier := createTestInstance()
		delivery := mockdelivery.CreateMockDelivery("test", "123", ([]byte)("[]"), true)
		go func() {
			supplier <- delivery
		}()
		_, _, err := instance.GetDelivery()
		if err == nil || !strings.HasPrefix(err.Error(), "Could not unmarshal json") {
			t.Error("Wrong exception", err)
		}
	})
	t.Run("Returns the correct metadata and log data", func(t *testing.T) {
		instance, _, supplier := createTestInstance()
		metadata := interfaces.LogMetadata{
			"test key": "test value",
		}
		message := rabbitMessage{
			metadata,
			dummydata.Row1,
		}
		json, _ := json.Marshal(message)
		delivery := mockdelivery.CreateMockDelivery("test", "123", json, true)
		go func() {
			supplier <- delivery
		}()
		receivedMetadata, receivedLogRow, err := instance.GetDelivery()
		if err != nil {
			t.Fail()
		}
		if !reflect.DeepEqual(metadata, receivedMetadata) {
			t.Error("Received wrong metadata", receivedMetadata)
		}
		if !reflect.DeepEqual(dummydata.Row1, receivedLogRow) {
			t.Error("Received wrong metadata", receivedLogRow)
		}
	})
}

func TestClose(t *testing.T) {
	t.Run("Closes the connecton", func(t *testing.T) {
		instance, channel, _ := createTestInstance()
		instance.Close()
		if !channel.isClosed {
			t.Error("Channel was not closed")
		}
	})
}
