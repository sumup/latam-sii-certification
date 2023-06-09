package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sumup/sii-certification/internal/adapters"
	"github.com/sumup/sii-certification/internal/appconfig"
	"github.com/sumup/sii-certification/internal/batches"
	"github.com/sumup/sii-certification/internal/entities"
)

func Auth(ctx context.Context, s *batches.Gateway) (string, error) {
	seed, err := s.GetSeed(ctx)
	if err != nil {
		return "", fmt.Errorf("errGetSeedReturnsError")
	}
	token, err := s.GetToken(ctx, seed)
	if err != nil {
		return "", fmt.Errorf("errGetTokenReturnsError")
	}

	return token, nil
}
func main() {
	appConfig := appconfig.FromEnv()
	httpClient := http.DefaultClient
	httpAdapter := adapters.NewAdapter(httpClient)
	taxAuthorityGateway := batches.NewTaxAuthorityGateway(httpAdapter, appConfig.TaxAuthorityChile)
	ctx := context.TODO()
	token, err := Auth(ctx, taxAuthorityGateway)
	if err != nil {
		fmt.Errorf("Authentication error")
		return
	}
	//fmt.Println(token)
	batchesMatrix := GenerateBatches()
	taxAuthorityGateway.SendMany(ctx, token, batchesMatrix)
}

func GenerateBatches() [][]entities.Batch {
	batchesMatrix := make([][]entities.Batch, 0)
	// document 48 with amount distinct to zero
	batchesMatrix = append(batchesMatrix, []entities.Batch{
		{
			VatID:           "96965568-3", // valid rut
			Day:             "2023-04-12",
			DocumentType:    "48",
			Channel:         "1", // 1 for cnp
			Amount:          10000,
			NTransactions:   2,
			ExternalTrackID: "0",
		},
		{
			VatID:           "96978044-5", // valid rut
			Day:             "2023-04-10",
			DocumentType:    "48",
			Channel:         "1", // 1 for cnp
			Amount:          2042,
			NTransactions:   2,
			ExternalTrackID: "0",
		},
		{
			VatID:           "65537690-9", // valid rut
			Day:             "2023-04-10",
			DocumentType:    "48",
			Channel:         "0", // 0 for cnp
			Amount:          500,
			NTransactions:   2,
			ExternalTrackID: "0",
		},
		{
			VatID:           "39020493-0", // valid rut
			Day:             "2023-04-10",
			DocumentType:    "48",
			Channel:         "0", // 0 for cnp
			Amount:          500,
			NTransactions:   2,
			ExternalTrackID: "0",
		},
	})

	// document 48 with amount zero
	batchesMatrix = append(batchesMatrix, []entities.Batch{
		{
			VatID:           "69610726-2", // valid rut
			Day:             "2023-04-10",
			DocumentType:    "48",
			Channel:         "1", // 1 for cnp
			Amount:          0,
			NTransactions:   1,
			ExternalTrackID: "0",
		},
		{
			VatID:           "28474336-9", // valid rut
			Day:             "2023-04-10",
			DocumentType:    "48",
			Channel:         "1", // 1 for cnp
			Amount:          0,
			NTransactions:   1,
			ExternalTrackID: "0",
		},
		{
			VatID:           "85445473-0", // valid rut
			Day:             "2023-04-10",
			DocumentType:    "48",
			Channel:         "0", // 0 for cnp
			Amount:          0,
			NTransactions:   1,
			ExternalTrackID: "0",
		},
	})

	// document 33 (bill)
	batchesMatrix = append(batchesMatrix, []entities.Batch{
		{
			VatID:           "94690304-3", // valid rut
			Day:             "2023-04-10",
			DocumentType:    "33",
			Channel:         "1", // 1 for cnp
			Amount:          10000,
			NTransactions:   2,
			ExternalTrackID: "0",
		},
	})

	// document 99 (no sell)
	batchesMatrix = append(batchesMatrix, []entities.Batch{
		{
			VatID:           "99040414-3", // valid rut
			Day:             "2023-04-10",
			DocumentType:    "99",
			Channel:         "1", // 1 for cnp
			Amount:          10000,
			NTransactions:   2,
			ExternalTrackID: "0",
		},
	})

	// document 00 (unknown)
	batchesMatrix = append(batchesMatrix, []entities.Batch{
		{
			VatID:           "61116878-0", // valid rut
			Day:             "2023-04-10",
			DocumentType:    "00",
			Channel:         "1", // 1 for cnp
			Amount:          2000,
			NTransactions:   2,
			ExternalTrackID: "0",
		},
	})
	return batchesMatrix
}
