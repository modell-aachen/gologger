package mockrpcchannel

import (
	"errors"
	"github.com/streadway/amqp"
)

type Published struct {
	Exchange  string
	Key       string
	Mandatory bool
	Immediate bool
	Msg       amqp.Publishing
}

type MockRpcChannel struct {
	name        string
	isClosed    bool
	Publishings chan *Published
}

func CreateMockRpcChannel(name string, isClosed bool) *MockRpcChannel {
	return &MockRpcChannel{
		name,
		isClosed,
		make(chan *Published),
	}
}

func (instance *MockRpcChannel) Close() error {
	instance.isClosed = true
	return nil
}
func (instance *MockRpcChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	if instance.isClosed {
		return errors.New("Channel has been closed")
	}
	published := Published{
		Exchange:  exchange,
		Key:       key,
		Mandatory: mandatory,
		Immediate: immediate,
		Msg:       msg,
	}
	instance.Publishings <- &published

	return nil
}

func (instance *MockRpcChannel) IsClosed() bool {
	return instance.isClosed
}
