package entities

import (
	"time"
)

type Transaction struct {
	TransactionID      string `gorm:"primarykey"`
	BatchID            *uint
	MerchantCode       string
	VatID              string
	Date               time.Time
	GrossTaxableAmount uint64
	GrossExemptAmount  uint64
	IsCNP              int
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Batch struct {
	ID              uint
	MerchantCode    string `gorm:"index:b_unique,unique"`
	Day             string `gorm:"index:b_unique,unique"`
	HasTaxes        bool   `gorm:"index:b_unique,unique"`
	IsCNP           bool   `gorm:"index:b_unique,unique"`
	Amount          uint64
	VatID           string
	Status          int
	LastResponse    string `gorm:"type:text"`
	ExternalTrackID string
	Transactions    []Transaction
	NTransactions   int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (u *Batch) AddTransactionAndTaxableAmount(t Transaction) {
	u.Amount += t.GrossTaxableAmount
	u.Transactions = append(u.Transactions, t)
}

func (u *Batch) AddTransactionAndExemptAmount(t Transaction) {
	u.Amount += t.GrossExemptAmount
	u.Transactions = append(u.Transactions, t)
}

func (u *Batch) UpdateStatus(status int) {
	u.Status = status
}

func NewBatch(merchantID string, vatID string, day string, hasTaxes bool, isCNP bool) Batch {
	return Batch{
		MerchantCode: merchantID,
		VatID:        vatID,
		Day:          day,
		Status:       PendingStatus,
		HasTaxes:     hasTaxes,
		IsCNP:        isCNP,
	}
}

const (
	PendingStatus   int = 1
	FailedStatus    int = 2
	SentButRejected int = 100
	SentAndAccepted int = 200
)

type Merchant struct {
	Code          string `json:"merchant_code"`
	LegalTypeID   int64  `json:"legal_type_id"`
	LegalTypeDesc string `json:"legal_type_desc"`
	NationalID    string `json:"national_id"`
	TaxID         string `json:"tax_id"`
	VatID         string `json:"vat_id"`
	IsSoleTrader  bool   `json:"is_sole_trader"`
	CountryCode   string `json:"country_code"`
}

type StatusesOfDay struct {
	Day             string
	SentAndAccepted int
	SentButRejected int
	Pending         int
	Failed          int
	Total           int
}

type DWHTransaction struct {
	TransactionCode string `json:"transaction_code"`
	PaymentType     int    `json:"payment_type"`
}
