package utils

import (
	"encoding/xml"
	"testing"
	"time"

	testify "github.com/stretchr/testify/assert"
)

func TestCalculateLimitDatesForMonth(t *testing.T) {
	t.Parallel()

	givenYear := "2024"
	givenMonth := "02"

	expectedFirstDay := "2024-02-01"
	expectedFirstDayNextMonth := "2024-03-01"

	resolvedFirstDay, resolvedFirstDayNextMonth := CalculateLimitDatesForMonth(givenMonth, givenYear)

	if resolvedFirstDay != expectedFirstDay {
		t.Errorf("got %q, wanted %q", resolvedFirstDay, expectedFirstDay)
	}

	if resolvedFirstDayNextMonth != expectedFirstDayNextMonth {
		t.Errorf("got %q, wanted %q", resolvedFirstDayNextMonth, expectedFirstDayNextMonth)
	}
}

func TestCalculateLimitDatesForMonth_EndOfYear(t *testing.T) {
	t.Parallel()

	givenYear := "2021"
	givenMonth := "12"

	expectedFirstDay := "2021-12-01"
	expectedFirstDayNextMonth := "2022-01-01"

	resolvedFirstDay, resolvedFirstDayNextMonth := CalculateLimitDatesForMonth(givenMonth, givenYear)

	if resolvedFirstDay != expectedFirstDay {
		t.Errorf("got %q, wanted %q", resolvedFirstDay, expectedFirstDay)
	}

	if resolvedFirstDayNextMonth != expectedFirstDayNextMonth {
		t.Errorf("got %q, wanted %q", resolvedFirstDayNextMonth, expectedFirstDayNextMonth)
	}
}

func TestLastDayOfMonth(t *testing.T) {
	assert := testify.New(t)

	lastDay := LastDayOfMonth("01", "2024")

	assert.Equal("31-01-2024", lastDay)

	lastDay = LastDayOfMonth("02", "2024")

	assert.Equal("29-02-2024", lastDay)

	lastDay = LastDayOfMonth("02", "2025")

	assert.Equal("28-02-2025", lastDay)
}

func TestLimitDatesForDay(t *testing.T) {
	assert := testify.New(t)

	chileLocation, _ := time.LoadLocation("America/Santiago")

	startTime, endTime := LimitTimesForDay("CL", "2022", "08", "08")

	// assert.Equal("2022-08-08", startTime)
	// assert.Equal("2022-08-09", endTime)
	assert.Equal(time.Date(2022, time.August, 7, 20, 0, 0, 0, chileLocation), startTime)
	assert.Equal(time.Date(2022, time.August, 8, 20, 0, 0, 0, chileLocation), endTime)

	startTime, endTime = LimitTimesForDay("CL", "2022", "08", "31")

	assert.Equal(time.Date(2022, time.August, 30, 20, 0, 0, 0, chileLocation), startTime)
	assert.Equal(time.Date(2022, time.August, 31, 20, 0, 0, 0, chileLocation), endTime)

	startTime, endTime = LimitTimesForDay("CL", "2024", "02", "28")

	assert.Equal(time.Date(2024, time.February, 27, 21, 0, 0, 0, chileLocation), startTime)
	assert.Equal(time.Date(2024, time.February, 28, 21, 0, 0, 0, chileLocation), endTime)

	startTime, endTime = LimitTimesForDay("CL", "2024", "03", "01")

	assert.Equal(time.Date(2024, time.February, 29, 21, 0, 0, 0, chileLocation), startTime)
	assert.Equal(time.Date(2024, time.March, 01, 21, 0, 0, 0, chileLocation), endTime)

}

func TestMergeIDs(t *testing.T) {
	assert := testify.New(t)

	var ids = make([]*int64, 0)
	merged := MergeIDs(ids)

	assert.Equal("", merged)

	var value1 int64 = 1
	var value2 int64 = 2
	var value3 int64 = 3

	ids = []*int64{&value1, &value2, &value3}
	merged = MergeIDs(ids)

	assert.Equal("1,2,3", merged)

	var value4 int64 = 6

	ids = []*int64{&value4}
	merged = MergeIDs(ids)

	assert.Equal("6", merged)
}

type TokenObject struct {
	XMLName xml.Name        `xml:"getToken"`
	Item    TokenObjectItem `xml:",omitempty"`
}

type TokenObjectItem struct {
	XMLName xml.Name `xml:"item"`
	Seed    string   `xml:"Semilla,omitempty"`
}

func (o *TokenObject) FillSeed(seed string) {
	o.Item.Seed = seed
}

func TestCanonalize(t *testing.T) {
	assert := testify.New(t)

	var seed = "000002248802"
	var tokenObject = TokenObject{}
	tokenObject.FillSeed(seed)

	canon, _ := CanonicalizeXML(tokenObject)

	assert.Equal("<getToken><item><Semilla>000002248802</Semilla></item></getToken>", string(canon))
}

func TestHash64(t *testing.T) {
	assert := testify.New(t)

	payload := []byte("<SignedInfo><CanonicalizationMethod Algorithm=\"http://www.w3.org/TR/2001/REC-xml-c14n-20010315\"/><SignatureMethod Algorithm=\"http://www.w3.org/2000/09/xmldsig#rsa-sha1\"/><Reference URI=\"\"><Transforms><Transform Algorithm=\"http://www.w3.org/2000/09/xmldsig#enveloped-signature\"/></Transforms><DigestMethod Algorithm=\"http://www.w3.org/2000/09/xmldsig#sha1\"/><DigestValue>kZvDbarenZxZPbWY7gNLxOan/NI=</DigestValue></Reference></SignedInfo>")

	assert.Equal("1cxbp15zpBlhnEgFpBqWPJ73Has=", Sha1HashBase64(payload))
}
