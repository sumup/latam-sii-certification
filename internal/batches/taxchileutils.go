package batches

import (
	"math"
)

// ChileanRound takes input amount float
// returns rounded amount to the closest integer.
func ChileanRound(amount float64) int {
	return int(math.Round(amount))
}

func documentType(hasTaxes bool) string {
	if hasTaxes {
		return "48" // 00: Has IVA
	}

	return "99" // Exempt from IVA
}

func channelType(isCNP bool) string {
	if isCNP {
		return "1"
	}

	return "0"
}

func trackIDOrDefault(id string) string {
	if id != "" {
		return id
	}

	return "0"
}
