package stripe

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

type StripePaymentGateway struct {
	store         TransactionStore
	validator     Validator
	authenticator Authenticator
}

func NewStripePaymentGateway(
	store TransactionStore,
	validator Validator,
	authenticator Authenticator,
) *StripePaymentGateway {
	return &StripePaymentGateway{
		store:         store,
		validator:     validator,
		authenticator: authenticator,
	}
}

func (spg *StripePaymentGateway) ProcessPayment(ctx context.Context, details gateway.PaymentDetails) (*gateway.TransactionStatus, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if err := spg.validator.Validate(details); err != nil {
		return nil, err
	}

	if !spg.authenticator.Authenticate(details.CardNumber) {
		return nil, errors.New("stripe: authentication failed: invalid card")
	}

	saveTransaction := &SaveTransaction{
		Amount:   details.Amount,
		Currency: details.Currency,
	}

	savedTransaction, err := spg.store.Save(ctx, saveTransaction)

	if err != nil {
		return nil, err
	}

	return savedTransaction, nil
}

func (spg *StripePaymentGateway) Refund(ctx context.Context, details gateway.RefundDetails) (*gateway.TransactionStatus, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	transaction, err := spg.store.Get(ctx, details.TransactionID)

	if err != nil {
		return nil, err
	}

	if transaction.Status != gateway.StatusCompleted {
		return nil, errors.New("stripe: cannot refund non-completed transaction")
	}

	if details.Amount > transaction.Amount {
		return nil, errors.New("stripe: refund amount exceeds transaction amount")
	}

	updateTransaction := &UpdateTransaction{
		TransactionID: details.TransactionID,
		Status:        gateway.StatusRefund,
	}

	updatedTransaction, err := spg.store.Update(ctx, *updateTransaction)

	if err != nil {
		return nil, err
	}

	return updatedTransaction, nil
}

func (spg *StripePaymentGateway) GetStatus(ctx context.Context, transactionID string) (*gateway.TransactionStatus, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	transaction, err := spg.store.Get(ctx, transactionID)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}
