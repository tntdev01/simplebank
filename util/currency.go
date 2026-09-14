package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

// IsValidCurrency checks if the given currency is in the list of valid currencies.
func IsValidCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD:
		return true
	}
	return false
}
