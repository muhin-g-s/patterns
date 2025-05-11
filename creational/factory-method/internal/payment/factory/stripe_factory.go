package factory

import (
	"factory-method/internal/payment/gateway"
	"factory-method/internal/payment/gateway/stripe"
)

type StripeGatewayFactory struct{}

func NewStripeGatewayFactory() *StripeGatewayFactory {
	return &StripeGatewayFactory{}
}

func (*StripeGatewayFactory) GetPaymentGateway() gateway.PaymentGateway {
	validator := stripe.NewDefaultValidator()
	authenticator := stripe.NewSimpleCardAuthenticator()
	store := stripe.NewTransactionStore()

	gateway := stripe.NewStripePaymentGateway(store, validator, authenticator)

	return gateway
}
