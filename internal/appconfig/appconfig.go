package appconfig

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Env               string
	DriveBaseFolder   string
	TaxAuthorityChile TaxAutChile
	Flags
}

type TaxAutChile struct {
	BaseURL                  string
	AuthURL                  string
	SendTransactionsEndpoint string
	GetSeedEndpoint          string
	GetTokenEndpoint         string
	SumUpRutBase             string
	SumUpRutVerifier         string
	SumUpCertificate         string
	SumUpCertificateModulus  string
	SumUpCertificateExponent string
}

type Flags struct {
	CheckTaxToggle *bool
}

func FromEnv() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	// reading and parsing flag
	flag.CommandLine = flag.NewFlagSet("check_tax_toggle", flag.ExitOnError)
	checkTaxToggle := flag.Bool("check_tax_toggle", false, "check in DWH if tax is enabled")
	flag.Parse()

	return Config{
		os.Getenv("ENV"),
		os.Getenv("DRIVE_BASE_FOLDER"),
		TaxAutChile{
			BaseURL:                  os.Getenv("TAX_AUTHORITY_CHILE_URL"),
			AuthURL:                  os.Getenv("TAX_AUTHORITY_CHILE_AUTH_URL"),
			SendTransactionsEndpoint: "/comprobanteboletaservice",
			GetSeedEndpoint:          "/CrSeed.jws",
			GetTokenEndpoint:         "/GetTokenFromSeed.jws",
		},
		Flags{
			CheckTaxToggle: checkTaxToggle,
		},
	}
}
