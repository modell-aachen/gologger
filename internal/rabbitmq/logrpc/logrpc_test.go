package logrpc

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/modell-aachen/gologger/internal/interfaces"
	"github.com/modell-aachen/gologger/internal/mocks/deliverysupplier"
	"github.com/modell-aachen/gologger/internal/mocks/mockdelivery"
	"github.com/modell-aachen/gologger/internal/mocks/mockrpcchannel"
)

func createTestInstance() (*RpcConnection, *mockrpcchannel.MockRpcChannel, deliverysupplier.DeliverySupplier) {
	channel := mockrpcchannel.CreateMockRpcChannel("test", false)
	supplier := deliverysupplier.CreateDeliverySupplier()
	instance := CreateRpcConnection(channel, supplier)

	return instance, channel, supplier
}

func TestGetRequest(t *testing.T) {
	t.Run("Notices, when it can not decode delivery", func(t *testing.T) {
		instance, _, supplier := createTestInstance()
		go func() {
			delivery := mockdelivery.CreateMockDelivery("mock", "123", ([]byte)("[]"), true)
			supplier <- delivery
		}()
		_, _, _, _, _, err := instance.GetRequest()
		if err == nil || !strings.HasPrefix(err.Error(), "Could not unmarshal json") {
			t.Error("Did not get expected error", err)
		}
	})
	t.Run("Decodes the request correctly", func(t *testing.T) {
		instance, _, supplier := createTestInstance()
		request := RpcLogRequest{
			[]interfaces.LevelString{"level 1", "level 2"},
			time.Unix(1234, 0),
			time.Unix(9876, 0),
			interfaces.SourceString("mock source"),
		}
		go func() {
			json, err := json.Marshal(request)
			if err != nil {
				t.Fatal()
			}
			delivery := mockdelivery.CreateMockDelivery("mock", "123", json, true)
			supplier <- delivery
		}()
		_, startTime, endTime, source, levels, err := instance.GetRequest()
		if err != nil {
			t.Fail()
		}
		if !startTime.Equal(request.Start_time) {
			t.Error("Start time wrong")
		}
		if !endTime.Equal(request.End_time) {
			t.Error("End time wrong")
		}
		if source != request.Source {
			t.Error("Source wrong")
		}
		if len(levels) != 2 || levels[0] != request.Levels[0] || levels[1] != request.Levels[1] {
			t.Error("Levels wrong")
		}
	})
}

func TestClose(t *testing.T) {
	t.Run("Closes the connecton", func(t *testing.T) {
		instance, channel, _ := createTestInstance()
		instance.Close()
		if !channel.IsClosed() {
			t.Error("Channel was not closed")
		}
	})
}
