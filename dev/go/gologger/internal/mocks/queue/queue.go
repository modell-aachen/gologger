package queue

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/modell-aachen/gologger/internal/interfaces"
	"github.com/modell-aachen/gologger/internal/rabbitmq/logrpc"

	"github.com/modell-aachen/gologger/internal/mocks/mockdelivery"
	"github.com/modell-aachen/gologger/internal/mocks/mockrpcchannel"
)

func CreateQueueMock() *QueueMock {
	rpcChannels := make(map[string]*mockrpcchannel.MockRpcChannel)
	rpcDeliveryChannel := make(chan *mockdelivery.MockDelivery)
	MockDeliveryChannel := make(chan *MockDelivery)
	mock := QueueMock{
		rpcChannels,
		rpcDeliveryChannel,
		MockDeliveryChannel,
		[]*logReceiverMock{},
		false,
	}

	return &mock
}

type QueueMock struct {
	rpcChannels         map[string]*mockrpcchannel.MockRpcChannel
	RpcDeliveries       chan *mockdelivery.MockDelivery
	MockDeliveryChannel chan *MockDelivery

	receivers []*logReceiverMock

	isClosed bool
}

var _ interfaces.QueueInstance = &QueueMock{}

type MockDelivery struct {
	Metadata    interfaces.LogMetadata
	LogRow      interfaces.LogRow
	ReportError error
}

type logReceiverMock struct {
	deliveries *chan *MockDelivery
	isClosed   bool
}

var _ interfaces.LogReceiver = (*logReceiverMock)(nil)

func (instance *logReceiverMock) GetDelivery() (metadata interfaces.LogMetadata, logRow interfaces.LogRow, err error) {
	delivery := <-*instance.deliveries
	return delivery.Metadata, delivery.LogRow, delivery.ReportError
}

func (instance *logReceiverMock) Close() {
	instance.isClosed = true
}

func (instance *logReceiverMock) IsClosed() bool {
	return instance.isClosed
}

func (instance *QueueMock) Close() (err error) {
	instance.isClosed = true
	return nil
}

func (instance *QueueMock) IsClosed() bool {
	return instance.isClosed
}

func (instance *QueueMock) GetReceiver(name string) (logReceiver interfaces.LogReceiver, err error) {
	if instance.isClosed {
		return nil, errors.New("Mock connection closed")
	}
	receiver := &logReceiverMock{
		&instance.MockDeliveryChannel,
		false,
	}
	instance.receivers = append(instance.receivers, receiver)
	return receiver, nil
}

func (instance *QueueMock) GetDelivery() (delivery interfaces.Delivery) {
	delivery = <-instance.RpcDeliveries
	return delivery
}

func (instance *QueueMock) GetRpcChannel(name string) *mockrpcchannel.MockRpcChannel {
	channel, ok := instance.rpcChannels[name]
	if !ok {
		channel = mockrpcchannel.CreateMockRpcChannel(name, false)
		instance.rpcChannels[name] = channel
	}
	return channel
}

func (instance *QueueMock) MockDelivery(metadata interfaces.LogMetadata, logRow interfaces.LogRow, reportError error) {
	instance.MockDeliveryChannel <- &MockDelivery{
		metadata,
		logRow,
		reportError,
	}
}

func (instance *QueueMock) GetRpcReceiver() (logRpc interfaces.LogRpc, err error) {
	if instance.isClosed {
		return nil, errors.New("Connection has been closed")
	}

	channelName := "rpc_channel"
	channel := instance.GetRpcChannel(channelName)

	logRpc = logrpc.CreateRpcConnection(
		channel,
		instance,
	)
	return logRpc, nil
}

func (instance *QueueMock) MockMalformedRpcRequest(correlationId string, canAck bool) {
	rpcDelivery := mockdelivery.CreateMockDelivery(
		"rpc_channel",
		correlationId,
		([]byte)("malformed!"),
		canAck,
	)
	instance.RpcDeliveries <- rpcDelivery
}

func (instance *QueueMock) MockRpcRequest(correlationId string, canAck bool, startTime time.Time, endTime time.Time, source interfaces.SourceString, levels []interfaces.LevelString) {
	rpcLogRequest := logrpc.RpcLogRequest{
		levels,
		startTime,
		endTime,
		source,
	}
	body, err := json.Marshal(rpcLogRequest)
	if err != nil {
		log.Fatal("Could not stringify request", err)
	}

	rpcDelivery := mockdelivery.CreateMockDelivery(
		"rpc_channel",
		correlationId,
		body,
		canAck,
	)
	instance.RpcDeliveries <- rpcDelivery
}

func (instance *QueueMock) GetReceivers() (receivers []*logReceiverMock) {
	return instance.receivers
}

var _ interfaces.QueueInstance = (*QueueMock)(nil)
