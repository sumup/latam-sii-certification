package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/now"
	"github.com/rs/zerolog/log"
)

func LimitTimesForDay(countryCode, yearAsString, monthAsString, dayAsString string) (time.Time, time.Time) {
	loc := Location(countryCode)

	startDate, _ := time.Parse("2006-01-02", yearAsString+"-"+monthAsString+"-"+dayAsString)
	startDate = startDate.In(loc)

	endDate := startDate.AddDate(0, 0, 1)

	return startDate, endDate
}

func Location(countryCode string) *time.Location {
	fallbackLocation, _ := time.LoadLocation("America/Santiago")

	timezones := map[string]string{
		"PE": "America/Lima",
		"CL": "America/Santiago",
	}
	timezone, ok := timezones[countryCode]

	if !ok {
		log.Error().Msg("timezone not configured for country: " + countryCode)

		return fallbackLocation
	}

	location, err := time.LoadLocation(timezone)

	if err != nil {
		log.Error().Msg("cannot load location for timezone: " + timezone)

		return fallbackLocation
	}

	return location
}

func CalculateLimitDatesForMonth(monthAsString, yearAsString string) (string, string) {
	year, _ := strconv.Atoi(yearAsString)
	month, _ := strconv.Atoi(monthAsString)
	t := time.Date(year, time.Month(month), 15, 12, 30, 30, 123456789, time.UTC)
	firstDay := now.With(t).BeginningOfMonth()
	lastDay := now.With(t).EndOfMonth()
	firstDayNextMonth := lastDay.AddDate(0, 0, 1)

	return firstDay.Format("2006-01-02"), firstDayNextMonth.Format("2006-01-02")
}

func LastDayOfMonth(monthAsString, yearAsString string) string {
	year, _ := strconv.Atoi(yearAsString)
	month, _ := strconv.Atoi(monthAsString)
	t := time.Date(year, time.Month(month), 15, 12, 30, 30, 123456789, time.UTC)
	lastDay := now.With(t).EndOfMonth()

	return lastDay.Format("02-01-2006")
}

func CountryName(countryID int) (string, error) {
	countryNames := map[int]string{
		1084: "Peru",
	}
	countryName, ok := countryNames[countryID]

	if !ok {
		return "", ErrCountryNotConfigured
	}

	return countryName, nil
}

func CountryTaxMultiplier(countryID int) (float64, error) {
	taxMultipliers := map[int]float64{
		1084: 0.1525,
	}
	taxMultiplier, ok := taxMultipliers[countryID]

	if !ok {
		return 0.0, ErrCountryNotConfigured
	}

	return taxMultiplier, nil
}

func MergeIDs(ids []*int64) string {
	if len(ids) == 0 {
		return ""
	}

	var idsString = make([]string, len(ids))

	for index, value := range ids {
		idsString[index] = strconv.FormatInt(*value, 10)
	}

	return strings.Join(idsString, ",")
}

func CanonicalizeXML(object interface{}) ([]byte, error) {
	marshal, err := xml.Marshal(object)
	if err != nil {
		return []byte(""), errCantCanonicalizeXML
	}

	return marshal, nil
}

func Sha1HashBase64(payload []byte) string {
	h := sha1.New()
	h.Write(payload)
	encoded := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return encoded
}
