package forex

import (
	"errors"
	"os"
	"slices"
	"time"

	"github.com/mattevans/dinero"
)

type forexClient struct {
	baseCurrency     string
	targetCurrencies []string
	client           *dinero.Client
}

func New(baseCurrency string, targetCurrencies []string, cacheDuration time.Duration) (*forexClient, error) {
	auth_token, ok := os.LookupEnv("OPEN_EXCHANGE_APP_ID")
	if !ok {
		return &forexClient{}, ErrNoAuthToken
	}
	client := dinero.NewClient(
		auth_token,
		baseCurrency,
		cacheDuration,
	)

	availableCurrencies, err := client.Currencies.List()
	if err != nil {
		return &forexClient{}, errors.Join(ErrOpenExchangeAPIError, err)
	}

	availableCurrCodes := make([]string, 0, len(availableCurrencies))
	for _, currency := range availableCurrencies {
		availableCurrCodes = append(availableCurrCodes, currency.Code)
	}

	if !slices.Contains(availableCurrCodes, baseCurrency) {
		return &forexClient{}, &UnsupportedCurrenciesError{Currencies: []string{baseCurrency}}
	}

	diff := sliceDifference(targetCurrencies, availableCurrCodes)
	if len(diff) != 0 {
		return &forexClient{}, &UnsupportedCurrenciesError{Currencies: diff}
	}

	return &forexClient{
		baseCurrency:     baseCurrency,
		targetCurrencies: targetCurrencies,
		client:           client,
	}, nil
}

func (f *forexClient) GetRates() (map[string]float64, error) {
	result := make(map[string]float64)

	response, err := f.client.Rates.List()
	if err != nil {
		return result, err
	}

	for code, rate := range response.Rates {
		if !slices.Contains(f.targetCurrencies, code) {
			continue
		}

		result[code] = rate
	}

	return result, nil
}

// func sliceIntersection[T comparable](slice1, slice2 []T) []T {
// 	set := make(map[T]struct{})
// 	var result []T
//
// 	for _, item := range slice1 {
// 		set[item] = struct{}{}
// 	}
//
// 	for _, item := range slice2 {
// 		if _, found := set[item]; found {
// 			result = append(result, item)
// 			delete(set, item)
// 		}
// 	}
//
// 	return result
// }

func sliceDifference[T comparable](slice1, slice2 []T) []T {
	set := make(map[T]struct{})
	var result []T

	// Add all elements from slice2 to the map
	for _, item := range slice2 {
		set[item] = struct{}{}
	}

	// Check if elements from slice1 are not in the map
	for _, item := range slice1 {
		if _, found := set[item]; !found {
			result = append(result, item)
		}
	}

	return result
}
