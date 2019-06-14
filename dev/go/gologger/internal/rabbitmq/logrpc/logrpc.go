package logrpc

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/modell-aachen/gologger/internal/interfaces"
	"github.com/modell-aachen/gologger/internal/rabbitmq/logrpcrequest"
)

type RpcConnection struct {
	channel    interfaces.RpcChannel
	deliveries interfaces.DeliverySupplier
}

func CreateRpcConnection(channel interfaces.RpcChannel, deliveries interfaces.DeliverySupplier) *RpcConnection {
	return &RpcConnection{
		channel,
		deliveries,
	}
}

type RpcLogRequest struct {
	Levels     []interfaces.LevelString
	Start_time time.Time
	End_time   time.Time
	Source     interfaces.SourceString
}

func (instance *RpcConnection) GetRequest() (request interfaces.LogRpcRequest, startTime time.Time, endTime time.Time, source interfaces.SourceString, levels []interfaces.LevelString, err error) {
	delivery := instance.deliveries.GetDelivery()

	decoded := RpcLogRequest{}
	err = json.Unmarshal(delivery.GetBody(), &decoded)
	if err != nil {
		return request, startTime, endTime, source, levels, errors.Wrapf(err, "Could not unmarshal json: %s", delivery.GetBody())
	}

	startTime = decoded.Start_time
	endTime = decoded.End_time
	source = decoded.Source
	levels = decoded.Levels

	request = logrpcrequest.CreateLogRpcRequest(instance.channel, delivery)
	return request, startTime, endTime, source, levels, nil
}

func (instance *RpcConnection) Close() error {
	return instance.channel.Close()
}
