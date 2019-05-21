package rabbitmq

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/modell-aachen/gologger/interfaces"

	"github.com/streadway/amqp"
)

type rpcConnection struct {
	channel    *amqp.Channel
	deliveries <-chan amqp.Delivery
}

func (c rpcConnection) ReplyToRequest(delivery interfaces.RpcDelivery, fields []interfaces.LogRow) error {
	field_json, err := json.Marshal(fields)
	if err != nil {
		return errors.Wrap(err, "Could not marshal fields")
	}
	err = c.channel.Publish(
		"",
		delivery.GetReplyTo(),
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: delivery.GetCorrelationId(),
			Body:          field_json,
		},
	)
	if err != nil {
		return errors.Wrap(err, "Could not send reply")
	}

	err = delivery.Ack(true)
	if err != nil {
		return errors.Wrap(err, "Could not ACK delivery")
	}

	return nil
}

type rpcLogRequest struct {
	Levels     []interfaces.LevelString
	Start_time time.Time
	End_time   time.Time
	Source     interfaces.SourceString
}

func (c rpcConnection) GetRequest() (delivery interfaces.RpcDelivery, startTime time.Time, endTime time.Time, source interfaces.SourceString, levels []interfaces.LevelString, err error) {
	rabbitDelivery := <-c.deliveries
	delivery = &amqpDeliveryWrapper{
		&rabbitDelivery,
	}

	decoded := rpcLogRequest{}
	err = json.Unmarshal(delivery.GetBody(), &decoded)
	if err != nil {
		return delivery, startTime, endTime, source, levels, errors.Wrapf(err, "Could not unmarshal json: %s", delivery.GetBody())
	}

	startTime = decoded.Start_time
	endTime = decoded.End_time
	source = decoded.Source
	levels = decoded.Levels
	return delivery, startTime, endTime, source, levels, nil
}

func (connection rpcConnection) Close() {
	connection.channel.Close()
}
