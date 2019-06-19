package logrpcrequest

import (
	"encoding/json"

	"github.com/modell-aachen/gologger/internal/interfaces"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type LogRpcRequest struct {
	channel  interfaces.RpcChannel
	delivery interfaces.Delivery
}

var _ interfaces.LogRpcRequest = (*LogRpcRequest)(nil)

func CreateLogRpcRequest(channel interfaces.RpcChannel, delivery interfaces.Delivery) *LogRpcRequest {
	return &LogRpcRequest{
		channel,
		delivery,
	}
}

func (instance *LogRpcRequest) Reply(fields []interfaces.LogRow) (err error) {
	field_json, err := json.Marshal(fields)
	if err != nil {
		return errors.Wrap(err, "Could not marshal fields")
	}
	err = instance.channel.Publish(
		"",
		instance.delivery.GetReplyTo(),
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: instance.delivery.GetCorrelationId(),
			Body:          field_json,
		},
	)
	if err != nil {
		return errors.Wrap(err, "Could not send reply")
	}

	err = instance.delivery.Ack(true)
	if err != nil {
		return errors.Wrap(err, "Could not ACK delivery")
	}

	return nil
}
