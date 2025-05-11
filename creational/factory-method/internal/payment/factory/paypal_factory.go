package factory

import (
	"factory-method/internal/payment/gateway"
	"factory-method/internal/payment/gateway/paypal"
)

type PaypalGatewayFactory struct{}

func NewPaypalGatewayFactory() *PaypalGatewayFactory {
	return &PaypalGatewayFactory{}
}

func (*PaypalGatewayFactory) GetPaymentGateway() gateway.PaymentGateway {
	validator := paypal.NewDefaultValidator()
	authenticator := paypal.NewSimpleCardAuthenticator()
	store := paypal.NewTransactionStore()

	gateway := paypal.NewPaypalPaymentGateway(store, validator, authenticator)

	return gateway
}
