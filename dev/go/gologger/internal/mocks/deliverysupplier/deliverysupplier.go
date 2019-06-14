package deliverysupplier

import (
	"github.com/modell-aachen/gologger/internal/interfaces"
)

type DeliverySupplier chan interfaces.Delivery

func (supplier DeliverySupplier) GetDelivery() interfaces.Delivery {
	return <-supplier
}

func CreateDeliverySupplier() DeliverySupplier {
	return (DeliverySupplier)(make(chan interfaces.Delivery))
}
