package factory

import "factory-method/internal/payment/gateway"

type PaymentGatewayFactory interface {
	GetPaymentGateway() *gateway.PaymentGateway
}
