package rating

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// https://www.tradingview.com/symbols/KRX-035420/technicals/?solution=43000614331
func Test_Rating(t *testing.T) {
	// Define the ratings and their corresponding weights
	delta := 0.0001
	weights := map[string]int{
		"sell":    -1,
		"neutral": 0,
		"buy":     1,
	}

	// Given ratings count
	sellCount := 15
	neutralCount := 9
	buyCount := 2

	totalCount := sellCount + neutralCount + buyCount

	avgSellCount := float64(sellCount) / float64(totalCount)
	avgNeutralCount := float64(neutralCount) / float64(totalCount)
	avgBuyCount := float64(buyCount) / float64(totalCount)

	// Calculate the overall rating result
	avgTotalWeightedValue := (avgSellCount * float64(weights["sell"])) +
		(avgNeutralCount * float64(weights["neutral"])) +
		(avgBuyCount * float64(weights["buy"]))

	assert.InDelta(t, -0.50, avgTotalWeightedValue, delta, "Unexpected value")

	// Output the overall rating result
	overallRating := ""

	switch {
	case avgTotalWeightedValue >= -1.0 && avgTotalWeightedValue < -0.5:
		overallRating = "Strong Sell"
	case avgTotalWeightedValue >= -0.5 && avgTotalWeightedValue < -0.1:
		overallRating = "Sell"
	case avgTotalWeightedValue >= -0.1 && avgTotalWeightedValue <= 0.1:
		overallRating = "Neutral"
	case avgTotalWeightedValue > 0.1 && avgTotalWeightedValue <= 0.5:
		overallRating = "Buy"
	case avgTotalWeightedValue > 0.5 && avgTotalWeightedValue <= 1.0:
		overallRating = "Strong Buy"
	}

	fmt.Printf("Overall Rating Result: %s\n", overallRating)
	assert.Equal(t, "Sell", overallRating)
}
