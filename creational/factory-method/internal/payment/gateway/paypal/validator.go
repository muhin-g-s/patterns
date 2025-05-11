package paypal

import (
	"errors"
	"factory-method/internal/payment/gateway"
)

type DefaultValidator struct{}

func NewDefaultValidator() *DefaultValidator {
	return &DefaultValidator{}
}

func (v *DefaultValidator) Validate(details gateway.PaymentDetails) error {
	if details.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if details.Currency == "" {
		return errors.New("currency is required")
	}
	if details.CardNumber == "" {
		return errors.New("card number is required")
	}
	return nil
}
