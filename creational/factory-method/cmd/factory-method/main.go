package main

import (
	"factory-method/internal/payment/factory"
	"factory-method/internal/payment/processor"
	"factory-method/pkg/api"
)

func main() {
	paypalFactory := factory.NewPaypalGatewayFactory()
	stripeFactory := factory.NewStripeGatewayFactory()

	paypalProcessor := processor.NewProcessor(paypalFactory)
	stripeProcessor := processor.NewProcessor(stripeFactory)
	_ = api.NewHandler(paypalProcessor, stripeProcessor)
}
