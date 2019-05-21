package mocks

import "github.com/modell-aachen/gologger/interfaces"

func CreatePostgresMock() PostgresMock {
	mock := PostgresMock{}
	mock.Data = &mockData{}

	return mock
}

func CreateQueueMock() (interfaces.QueueInstance, *chan MockRpcRequest, *chan MockDelivery) {
	rpcChannel := make(chan MockRpcRequest)
	deliveryChannel := make(chan MockDelivery)
	mock := QueueMock{
		&rpcChannel,
		&deliveryChannel,
		false,
	}

	return &mock, &rpcChannel, &deliveryChannel
}
