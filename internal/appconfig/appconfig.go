package appconfig

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Env               string
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

	cert, err := getTaxAuthorityCertificate()
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading tax authority certificate")
	}

	return Config{
		os.Getenv("ENV"),
		TaxAutChile{
			BaseURL:                  os.Getenv("TAX_AUTHORITY_CHILE_URL"),
			AuthURL:                  os.Getenv("TAX_AUTHORITY_CHILE_AUTH_URL"),
			SendTransactionsEndpoint: "/comprobanteboletaservice",
			GetSeedEndpoint:          "/CrSeed.jws",
			GetTokenEndpoint:         "/GetTokenFromSeed.jws",
			SumUpRutBase:             os.Getenv("TAX_AUTHORITY_SUMUP_RUT_BASE"),
			SumUpRutVerifier:         os.Getenv("TAX_AUTHORITY_SUMUP_RUT_VD"),
			SumUpCertificate:         cert,
			SumUpCertificateModulus:  os.Getenv("TAX_AUTHORITY_SUMUP_CERTIFICATE_MODULUS"),
			SumUpCertificateExponent: os.Getenv("TAX_AUTHORITY_SUMUP_CERTIFICATE_EXPONENT"),
		},
		Flags{
			CheckTaxToggle: checkTaxToggle,
		},
	}
}

func getTaxAuthorityCertificate() (string, error) {
	dir, _ := os.Getwd()
	//fmt.Println(dir)
	file, err := ioutil.ReadFile(filepath.Clean(dir + "/tax-authority.pem"))

	if err != nil {
		return "", fmt.Errorf("errCannotReadTaxAuthorityCertificate")
	}
	//fmt.Println(string(file))

	return string(file), nil
}
