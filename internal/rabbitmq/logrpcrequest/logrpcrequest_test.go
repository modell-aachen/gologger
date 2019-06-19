package logrpcrequest

import (
	"strings"
	"testing"

	"github.com/modell-aachen/gologger/internal/mocks/mockdelivery"
	"github.com/modell-aachen/gologger/internal/mocks/mockrpcchannel"
)

func createRpcConnection(canAck bool) (*LogRpcRequest, *mockrpcchannel.MockRpcChannel, *mockdelivery.MockDelivery) {
	rpcChannel := mockrpcchannel.CreateMockRpcChannel("test", false)
	rpcDelivery := createDelivery(canAck)

	instance := CreateLogRpcRequest(rpcChannel, rpcDelivery)
	return instance, rpcChannel, rpcDelivery
}

func createDelivery(canAck bool) *mockdelivery.MockDelivery {
	return mockdelivery.CreateMockDelivery(
		"mock_rpc_delivery",
		"123",
		nil,
		canAck,
	)
}

func TestReply(t *testing.T) {
	t.Run("Notices, when it can not ack", func(t *testing.T) {
		instance, channel, _ := createRpcConnection(false)
		go func() {
			_ = <-channel.Publishings
		}()
		err := instance.Reply(nil)
		if !strings.HasPrefix(err.Error(), "Could not ACK delivery") {
			t.Error("Got wrong exception", err)
		}
	})
	t.Run("Acks deliveries", func(t *testing.T) {
		instance, channel, delivery := createRpcConnection(true)
		go func() {
			_ = <-channel.Publishings
		}()
		err := instance.Reply(nil)
		if err != nil {
			t.Error("Reply failed", err)
		}
		if !delivery.IsAcked() {
			t.Error("Delivery is not acked")
		}
	})
	t.Run("Notices, when it can not reply", func(t *testing.T) {
		instance, channel, _ := createRpcConnection(true)
		channel.Close()
		err := instance.Reply(nil)
		if !strings.HasPrefix(err.Error(), "Could not send reply") {
			t.Error("Got wrong exception", err)
		}
	})
	t.Run("Publishes to correct channel", func(t *testing.T) {
		instance, channel, _ := createRpcConnection(true)
		go func() {
			err := instance.Reply(nil)
			if err != nil {
				t.Error("Reply failed", err)
			}
		}()
		published := <-channel.Publishings
		if published.Key != "mock_rpc_delivery" {
			t.Error("Published to wrong channel", published.Key)
		}
	})
	t.Run("Publishes with correct CorrelationId", func(t *testing.T) {
		instance, channel, _ := createRpcConnection(true)
		go func() {
			err := instance.Reply(nil)
			if err != nil {
				t.Error("Reply failed", err)
			}
		}()
		published := <-channel.Publishings
		if published.Msg.CorrelationId != "123" {
			t.Error("Published to wrong correlation id", published.Msg.CorrelationId)
		}
	})
}
