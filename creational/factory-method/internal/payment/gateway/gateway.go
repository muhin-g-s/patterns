package gateway

import (
	"context"
	"time"
)

type PaymentDetails struct {
	Amount      float64
	Currency    string
	CardNumber  string
	CardHolder  string
	ExpiryDate  string
	CVV         string
	Description string
}

type RefundDetails struct {
	TransactionID string
	Amount        float64
	Reason        string
}

type TransactionStatus struct {
	TransactionID string
	Status        TransactionStatusType
	Amount        float64
	Currency      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ErrorMessage  string
}

type TransactionStatusType string

const (
	StatusPending   TransactionStatusType = "pending"
	StatusCompleted TransactionStatusType = "completed"
	StatusFailed    TransactionStatusType = "failed"
	StatusRefund    TransactionStatusType = "refund"
)

type PaymentGateway interface {
	ProcessPayment(ctx context.Context, details PaymentDetails) (*TransactionStatus, error)
	Refund(ctx context.Context, details RefundDetails) (*TransactionStatus, error)
	GetStatus(ctx context.Context, transactionID string) (*TransactionStatus, error)
}
