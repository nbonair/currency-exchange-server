package utils

func GetUniqueCurrencies(currencies []string) []string {
	set := make(map[string]struct{})
	for _, currency := range currencies {
		set[currency] = struct{}{}
	}
	var uniqueCurrencyList []string
	for currency := range set {
		uniqueCurrencyList = append(uniqueCurrencyList, currency)
	}
	return uniqueCurrencyList
}
