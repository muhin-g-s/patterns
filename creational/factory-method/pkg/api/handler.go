package api

import (
	"context"
	"errors"
	"fmt"

	"factory-method/internal/payment/gateway"
	"factory-method/internal/payment/processor"
)

type Handler struct {
	paypalProcessor *processor.Processor
	stripeProcessor *processor.Processor
}

func NewHandler(paypalProcessor *processor.Processor, stripeProcessor *processor.Processor) *Handler {
	return &Handler{
		paypalProcessor: paypalProcessor,
		stripeProcessor: stripeProcessor,
	}
}

type ProviderType string

const (
	StripeProvider ProviderType = "stripe"
	PaypalProvider ProviderType = "paypal"
)

type PaymentDetails struct {
	provider    ProviderType
	Amount      float64
	Currency    string
	CardNumber  string
	CardHolder  string
	ExpiryDate  string
	CVV         string
	Description string
}

type RefundDetails struct {
	provider      ProviderType
	TransactionID string
	Amount        float64
	Reason        string
}

type CheckStatusDetails struct {
	provider      ProviderType
	TransactionID string
}

type TransactionStatus struct {
	Status   TransactionStatusType
	Amount   float64
	Currency string
}

type TransactionStatusType string

const (
	StatusPending   TransactionStatusType = "pending"
	StatusCompleted TransactionStatusType = "completed"
	StatusFailed    TransactionStatusType = "failed"
	StatusRefund    TransactionStatusType = "refund"
)

var (
	ErrInvalidPaymentDetails = errors.New("invalid payment details")
	ErrInvalidRefundDetails  = errors.New("invalid refund details")
	ErrUnknownProvider       = errors.New("unknown payment gateway specified")
	ErrInternal              = errors.New("internal error")
)

func (h *Handler) MakePayment(ctx context.Context, details PaymentDetails) (*TransactionStatus, error) {
	if err := validatePaymentDetails(details); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidPaymentDetails, err)
	}

	p, err := h.resolveProcessor(details.provider)

	if err != nil {
		return nil, err
	}

	status, err := p.MakePayment(ctx, convertToProcessorPaymentDetails(details))
	if err != nil {
		return nil, ErrInternal
	}

	return convertFromProcessorStatus(status), nil
}

func (h *Handler) MakeRefund(ctx context.Context, details RefundDetails) (*TransactionStatus, error) {
	if details.TransactionID == "" || details.Amount <= 0 {
		return nil, fmt.Errorf("%w: transaction ID or amount is invalid", ErrInvalidRefundDetails)
	}

	p, err := h.resolveProcessor(details.provider)

	if err != nil {
		return nil, err
	}

	status, err := p.MakeRefund(ctx, convertToProcessorRefundDetails(details))
	if err != nil {
		return nil, ErrInternal
	}

	return convertFromProcessorStatus(status), nil
}

func (h *Handler) CheckStatus(ctx context.Context, details CheckStatusDetails) (*TransactionStatus, error) {
	if details.TransactionID == "" {
		return nil, fmt.Errorf("empty transaction ID provided")
	}

	p, err := h.resolveProcessor(details.provider)

	if err != nil {
		return nil, err
	}

	status, err := p.CheckStatus(ctx, details.TransactionID)
	if err != nil {
		return nil, ErrInternal
	}

	return convertFromProcessorStatus(status), nil
}

func (h *Handler) resolveProcessor(provider ProviderType) (*processor.Processor, error) {
	var p *processor.Processor

	switch provider {
	case StripeProvider:
		p = h.stripeProcessor
	case PaypalProvider:
		p = h.paypalProcessor
	default:
		return nil, fmt.Errorf("%w: %v", ErrUnknownProvider, provider)
	}

	return p, nil
}

func validatePaymentDetails(details PaymentDetails) error {
	if details.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if details.Currency == "" {
		return errors.New("currency is required")
	}
	if details.CardNumber == "" {
		return errors.New("card number is required")
	}
	if details.CardHolder == "" {
		return errors.New("card holder name is required")
	}
	if details.ExpiryDate == "" {
		return errors.New("expiry date is required")
	}
	if details.CVV == "" {
		return errors.New("CVV is required")
	}
	return nil
}

func convertToProcessorPaymentDetails(details PaymentDetails) gateway.PaymentDetails {
	return gateway.PaymentDetails{
		Amount:      details.Amount,
		Currency:    details.Currency,
		CardNumber:  details.CardNumber,
		CardHolder:  details.CardHolder,
		ExpiryDate:  details.ExpiryDate,
		CVV:         details.CVV,
		Description: details.Description,
	}
}

func convertToProcessorRefundDetails(details RefundDetails) gateway.RefundDetails {
	return gateway.RefundDetails{
		TransactionID: details.TransactionID,
		Amount:        details.Amount,
		Reason:        details.Reason,
	}
}

func convertFromProcessorStatus(status *gateway.TransactionStatus) *TransactionStatus {
	return &TransactionStatus{
		Status:   TransactionStatusType(status.Status),
		Amount:   status.Amount,
		Currency: status.Currency,
	}
}
