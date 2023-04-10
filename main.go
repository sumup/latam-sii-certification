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
	batchesMatrix := make([][]entities.Batch, 10)
	taxAuthorityGateway.SendMany(ctx, token, batchesMatrix)
}

func GenerateBatches(batches *[][]entities.Batch) {

}
