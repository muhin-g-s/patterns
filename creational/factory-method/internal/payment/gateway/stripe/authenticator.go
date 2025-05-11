package stripe

type SimpleCardAuthenticator struct {
	validCards map[string]bool
}

func NewSimpleCardAuthenticator() *SimpleCardAuthenticator {
	return &SimpleCardAuthenticator{
		validCards: map[string]bool{
			"4242424242424242": true,
			"5555555555554444": true,
		},
	}
}

func (a *SimpleCardAuthenticator) Authenticate(card string) bool {
	_, valid := a.validCards[card]
	return valid
}
