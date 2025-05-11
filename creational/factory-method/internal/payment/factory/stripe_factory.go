package factory

import (
	"factory-method/internal/payment/gateway"
	"factory-method/internal/payment/gateway/stripe"
)

func GetPaymentGateway() gateway.PaymentGateway {
	validator := stripe.NewDefaultValidator()
	authenticator := stripe.NewSimpleCardAuthenticator()
	store := stripe.NewTransactionStore()

	gateway := stripe.NewStripePaymentGateway(store, validator, authenticator)

	return gateway
}
