package util

const (
	USD = "USD"
	EUR = "EUR"
	NGN = "NGN"
	CAD = "CAD"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, NGN, CAD:
		return true
	}
	return false
}
