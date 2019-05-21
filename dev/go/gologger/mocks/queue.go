package mocks

import (
	"errors"
	"time"

	"github.com/modell-aachen/gologger/interfaces"
)

type QueueMock struct {
	RpcChannel   *chan MockRpcRequest
	LogChannel   *chan MockDelivery
	ReplyChannel *chan []interfaces.LogRow

	isClosed bool
}

type MockRpcRequest struct {
	delivery  interfaces.RpcDelivery
	startTime time.Time
	endTime   time.Time
	source    interfaces.SourceString
	levels    []interfaces.LevelString
}

type logRpcMock struct {
	requestChannel *chan MockRpcRequest
	replyChannel   *chan []interfaces.LogRow

	isClosed bool
}

type MockDelivery struct {
	metadata interfaces.LogMetadata
	logRow   interfaces.LogRow
}

type logReceiverMock struct {
	deliveries *chan MockDelivery
	isClosed   bool
}

func (instance *logReceiverMock) GetDelivery() (metadata interfaces.LogMetadata, logRow interfaces.LogRow, err error) {
	delivery := <-*instance.deliveries
	return delivery.metadata, delivery.logRow, nil
}

func (instance *logReceiverMock) Close() {
	instance.isClosed = true
}

func (instance *QueueMock) Close() (err error) {
	instance.isClosed = true
	return nil
}

func (instance *QueueMock) GetReceiver(name string) (logReceiver interfaces.LogReceiver, err error) {
	return &logReceiverMock{
		instance.LogChannel,
		false,
	}, nil
}

func (instance *QueueMock) GetRpcReceiver() (logRpc interfaces.LogRpc, err error) {
	if instance.isClosed {
		return nil, errors.New("Connection has been closed")
	}
	logRpc = &logRpcMock{
		instance.RpcChannel,
		false,
	}
	return logRpc, nil
}

func (instance *logRpcMock) GetRequest() (delivery interfaces.RpcDelivery, startTime time.Time, endTime time.Time, source interfaces.SourceString, levels []interfaces.LevelString, err error) {
	if instance.isClosed {
		return delivery, startTime, endTime, source, levels, errors.New("Rpc connection has been closed")
	}

	request := <-*instance.channel
	return request.delivery, request.startTime, request.endTime, request.source, request.levels, nil
}

func (instance *logRpcMock) Close() {
	instance.isClosed = true
}

func (instance *logRpcMock) ReplyToRequest(delivery interfaces.RpcDelivery, rows []interfaces.LogRow) (err error) {
	return nil
}
