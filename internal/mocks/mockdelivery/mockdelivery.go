package mockdelivery

import (
	"errors"

	"github.com/modell-aachen/gologger/internal/interfaces"
)

type MockDelivery struct {
	replyTo       string
	correlationId string
	body          []byte
	canAck        bool
	isAcked       bool
}

var _ interfaces.Delivery = (*MockDelivery)(nil)

func CreateMockDelivery(replyTo string, correlationId string, body []byte, canAck bool) *MockDelivery {
	return &MockDelivery{
		replyTo,
		correlationId,
		body,
		canAck,
		false,
	}
}

func (instance *MockDelivery) GetReplyTo() string {
	return instance.replyTo
}

func (instance *MockDelivery) GetCorrelationId() string {
	return instance.correlationId
}

func (instance *MockDelivery) GetBody() []byte {
	return instance.body
}

func (instance *MockDelivery) Ack(multiple bool) error {
	if instance.canAck {
		if instance.isAcked {
			return errors.New("Already acked")
		} else {
			instance.isAcked = true
			return nil
		}
	} else {
		return errors.New("Can not ack")
	}
}

func (instance *MockDelivery) IsAcked() bool {
	return instance.isAcked
}
