package paypal

type SimpleCardAuthenticator struct {
	validCards map[string]bool
}

func NewSimpleCardAuthenticator() *SimpleCardAuthenticator {
	return &SimpleCardAuthenticator{
		validCards: map[string]bool{
			"1": true,
			"2": true,
		},
	}
}

func (a *SimpleCardAuthenticator) Authenticate(card string) bool {
	_, valid := a.validCards[card]
	return valid
}
