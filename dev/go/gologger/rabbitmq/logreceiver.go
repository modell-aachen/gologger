package rabbitmq

import (
	"encoding/json"
	"github.com/pkg/errors"

	"github.com/modell-aachen/gologger/interfaces"

	"github.com/streadway/amqp"
)

type logConnection struct {
	channel    *amqp.Channel
	deliveries <-chan amqp.Delivery
}

func (c logConnection) Close() {
	c.channel.Close()
}

type rabbitMessage struct {
	Metadata interfaces.LogMetadata
	Log_data interfaces.LogRow
}

type rabbitMessagePlain struct {
	metadata map[string]interface{}
	//log_data interfaces.LogRow
}

func (c logConnection) GetDelivery() (logMetadata interfaces.LogMetadata, logRow interfaces.LogRow, err error) {
	delivery := <-c.deliveries
	err = delivery.Ack(true)
	if err != nil {
		return logMetadata, logRow, errors.Wrap(err, "Could not ACK delivery")
	}

	message := rabbitMessage{}
	err = json.Unmarshal(delivery.Body, &message)
	if err != nil {
		return logMetadata, logRow, errors.Wrapf(err, "Could not unmarshal json: %s", delivery.Body)
	}
	logMetadata = message.Metadata
	logRow = message.Log_data

	err = json.Unmarshal(delivery.Body, &message)

	return logMetadata, logRow, nil
}
