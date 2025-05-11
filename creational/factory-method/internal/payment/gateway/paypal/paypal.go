package paypal

import (
	"context"
	"errors"
	"factory-method/internal/payment/gateway"
)

type Validator interface {
	Validate(gateway.PaymentDetails) error
}

type SaveTransaction struct {
	Amount   float64
	Currency string
}

type UpdateTransaction struct {
	TransactionID string
	Status        gateway.TransactionStatusType
}

type TransactionStore interface {
	Save(ctx context.Context, saveTransaction *SaveTransaction) (*gateway.TransactionStatus, error)
	Get(ctx context.Context, id string) (*gateway.TransactionStatus, error)
	Update(ctx context.Context, updateTransaction UpdateTransaction) (*gateway.TransactionStatus, error)
}

type Authenticator interface {
	Authenticate(card string) bool
}

type PaypalPaymentGateway struct {
	store         TransactionStore
	validator     Validator
	authenticator Authenticator
}

func NewPaypalPaymentGateway(
	store TransactionStore,
	validator Validator,
	authenticator Authenticator,
) *PaypalPaymentGateway {
	return &PaypalPaymentGateway{
		store:         store,
		validator:     validator,
		authenticator: authenticator,
	}
}

func (ppg *PaypalPaymentGateway) ProcessPayment(ctx context.Context, details gateway.PaymentDetails) (*gateway.TransactionStatus, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if err := ppg.validator.Validate(details); err != nil {
		return nil, err
	}

	if !ppg.authenticator.Authenticate(details.CardNumber) {
		return nil, errors.New("paypal: authentication failed: invalid card")
	}

	saveTransaction := &SaveTransaction{
		Amount:   details.Amount,
		Currency: details.Currency,
	}

	savedTransaction, err := ppg.store.Save(ctx, saveTransaction)

	if err != nil {
		return nil, err
	}

	return savedTransaction, nil
}

func (ppg *PaypalPaymentGateway) Refund(ctx context.Context, details gateway.RefundDetails) (*gateway.TransactionStatus, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	transaction, err := ppg.store.Get(ctx, details.TransactionID)

	if err != nil {
		return nil, err
	}

	if transaction.Status != gateway.StatusCompleted {
		return nil, errors.New("paypal: cannot refund non-completed transaction")
	}

	if details.Amount > transaction.Amount {
		return nil, errors.New("paypal: refund amount exceeds transaction amount")
	}

	updateTransaction := &UpdateTransaction{
		TransactionID: details.TransactionID,
		Status:        gateway.StatusRefund,
	}

	updatedTransaction, err := ppg.store.Update(ctx, *updateTransaction)

	if err != nil {
		return nil, err
	}

	return updatedTransaction, nil
}

func (ppg *PaypalPaymentGateway) GetStatus(ctx context.Context, transactionID string) (*gateway.TransactionStatus, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	transaction, err := ppg.store.Get(ctx, transactionID)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}
