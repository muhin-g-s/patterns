package processor

import (
	"context"
	"factory-method/internal/payment/gateway"
	"fmt"
)

type Processor struct {
	gateway gateway.PaymentGateway
}

func NewProcessor(gateway gateway.PaymentGateway) *Processor {
	return &Processor{
		gateway: gateway,
	}
}

func (p *Processor) MakePayment(ctx context.Context, details gateway.PaymentDetails) (*gateway.TransactionStatus, error) {
	status, err := p.gateway.ProcessPayment(ctx, details)
	if err != nil {
		return nil, fmt.Errorf("gateway ProcessPayment failed: %w", err)
	}
	return status, nil
}

func (p *Processor) MakeRefund(ctx context.Context, details gateway.RefundDetails) (*gateway.TransactionStatus, error) {
	status, err := p.gateway.Refund(ctx, details)
	if err != nil {
		return nil, fmt.Errorf("gateway Refund failed: %w", err)
	}
	return status, nil
}

func (p *Processor) CheckStatus(ctx context.Context, transactionID string) (*gateway.TransactionStatus, error) {
	return p.gateway.GetStatus(ctx, transactionID)
}
