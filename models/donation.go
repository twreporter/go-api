package models

import (
	"time"
)

type PayByPrimeDonation struct {
	ID                       uint       `gorm:"primary_key" json:"id"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
	DeletedAt                *time.Time `json:"deleted_at"`
	Status                   int        `gorm:"not null" json:"status"`
	Msg                      string     `gorm:"type:varchar(100);not null" json:"msg"`
	RecTradeID               string     `gorm:"type:varchar(20);not null" json:"rec_trade_id"`
	BankTransactionID        string     `gorm:"type:varchar(50);not null" json:"bank_transaction_id"`
	AuthCode                 string     `gorm:"type:varchar(6);not null" json:"auth_code"`
	Amount                   uint       `gorm:"not null" json:"amount"`
	Currency                 string     `gorm:"type:varchar(3);default:'TWD';not null" json:"currency"`
	OrderNumber              string     `gorm:"type:varchar(50);not null" json:"order_number"`
	Acquirer                 string     `gorm:"type:varchar(50);not null" json:"acquirer"`
	TransactionTime          *time.Time `json:"transaction_time"`
	BankTransactionStartTime *time.Time `json:"bank_transaction_start_time"`
	BankTransactionEndTime   *time.Time `json:"bank_transaction_end_time"`
	BankResultCode           *string    `gorm:"type:varchar(50)" json:"bank_result_code"`
	BankResultMsg            *string    `gorm:"type:varchar(50)" json:"bank_result_msg"`
	PayMethod                string     `gorm:"type:ENUM('credit_card','line','apple','google','samsung');not null;index:idx_pay_by_prime_donations_cardholder_email_pay_method" json:"pay_method"`
	CardholderEmail          string     `gorm:"type:varchar(100);not null;index:idx_pay_by_prime_donations_cardholder_email_pay_method" json:"cardholder_email"`
	CardholderPhoneNumber    *string    `gorm:"type:varchar(20)" json:"cardholder_phone_number"`
	CardholderName           *string    `gorm:"type:varchar(30)" json:"cardholder_name"`
	CardholderZipCode        *string    `gorm:"type:varchar(10)" json:"carholder_zip_code"`
	CardholderAddress        *string    `gorm:"type:varchar(100)" json:"carholder_address"`
	CardholderNationalID     *string    `gorm:"type:varchar(20)" json:"carholder_national_id"`
	CardInfoBinCode          *string    `gorm:"type:varchar(6)" json:"card_info_bin_code"`
	CardInfoLastFour         *string    `gorm:"type:varchar(4)" json:"card_info_last_four"`
	CardInfoIssuer           *string    `gorm:"type:varchar(50)" json:"card_info_issuer"`
	CardInfoFunding          *uint      `gorm:"type:tinyint" json:"card_info_funding"`
	CardInfoType             *uint      `gorm:"type:tinyint" json:"card_info_type"`
	CardInfoLevel            *string    `gorm:"type:varchar(10)" json:"card_info_level"`
	CardInfoCountry          *string    `gorm:"type:varchar(30)" json:"card_info_country"`
	CardInfoCountryCode      *string    `gorm:"type:varchar(10)" json:"card_info_country_code"`
	CardInfoExpiryDate       *string    `gorm:"type:varchar(6);" json:"card_info_expiry_date"`
	Details                  string     `gorm:"type:varchar(50);not null" json:"details"`
	MerchantID               string     `gorm:"type:varchar(30);not null" json:"merchant_id"`
}

type PayByCardTokenDonation struct {
	ID                       uint       `gorm:"primary_key" json:"id"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
	DeletedAt                *time.Time `json:"deleted_at"`
	PeriodicID               uint       `gorm:"not null;index:idx_pay_by_card_token_donations_periodic_id" json:"periodic_id"`
	Status                   int        `gorm:"not null" json:"status"`
	Msg                      string     `gorm:"type:varchar(100);not null" json:"msg"`
	RecTradeID               string     `gorm:"type:varchar(20);not null" json:"rec_trade_id"`
	BankTransactionID        string     `gorm:"type:varchar(50);not null" json:"bank_transaction_id"`
	AuthCode                 string     `gorm:"type:varchar(6);not null" json:"auth_code"`
	Amount                   uint       `gorm:"not null;index:idx_pay_by_card_token_donations_amount" json:"amount"`
	Currency                 string     `gorm:"type:varchar(3);default:'TWD';not null" json:"currency"`
	OrderNumber              string     `gorm:"type:varchar(50);not null" json:"order_number"`
	Acquirer                 string     `gorm:"type:varchar(50);not null" json:"acquirer"`
	TransactionTime          *time.Time `json:"transaction_time"`
	BankTransactionStartTime *time.Time `json:"bank_transaction_start_time"`
	BankTransactionEndTime   *time.Time `json:"bank_transaction_end_time"`
	BankResultCode           string     `gorm:"type:varchar(50)" json:"bank_result_code"`
	BankResultMsg            string     `gorm:"type:varchar(50)" json:"bank_result_msg"`
	Details                  string     `gorm:"type:varchar(50);not null" json:"details"`
	MerchantID               string     `gorm:"type:varchar(30);not null" json:"merchant_id"`
}

type PayByOtherMethodDonation struct {
	ID         uint       `gorm:"primary_key" json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
	Email      string     `gorm:"type:varchar(100);not null" json:"email"`
	PayMethod  string     `gorm:"type:varchar(50);not null;index:idx_pay_by_other_donations_pay_method" json:"pay_method"`
	Amount     uint       `gorm:"type:int(10) unsigned;index:idx_pay_by_other_donations_amount" json:"amount"`
	Curreny    string     `gorm:"type:char(3);default:'TWD';not null" json:"currency"`
	Details    string     `gorm:"type:varchar(50);not null" json:"details"`
	MerchantID string     `gorm:"type:varchar(30);not null" json:"merchant_id"`
}

type PeriodicDonation struct {
	ID                    uint       `gorm:"primary_key" json:"id"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	DeletedAt             *time.Time `json:"deleted_at"`
	CardToken             string     `gorm:"type:varchar(64);not null" json:"card_token"`
	CardKey               string     `gorm:"type:varchar(64);not null" json:"card_key"`
	UserID                uint       `gorm:"type:int(10) unsigned;not null" json:"user_id"`
	Currency              string     `gorm:"type:varchar(3);default:'TWD';not null" json:"currency"`
	StartDate             time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;index:idx_periodic_donations_start_date" json:"start_date"`
	PaidTimes             uint       `gorm:"type:smallint;default:1;not null" json:"paid_times"`
	Amount                uint       `gorm:"type:int(10) unsigned;not null;index:idx_periodic_donations_amount" json:"amount"`
	LastSuccessDate       time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"last_success_date"`
	FailureTimes          uint       `gorm:"type:tinyint unsigned;default:0;not null" json:"failure_times"`
	IsStopped             bool       `gorm:"type:tinyint(1) unsigned;default:0;not null;index:idx_periodic_donations_is_stopped" json:"is_stopped"`
	CardholderEmail       string     `gorm:"type:varchar(100);not null" json:"cardholder_email"`
	CardholderPhoneNumber *string    `gorm:"type:varchar(20)" json:"cardholder_phone_number"`
	CardholderName        *string    `gorm:"type:varchar(30)" json:"cardholder_name"`
	CardholderZipCode     *string    `gorm:"type:varchar(10)" json:"carholder_zip_code"`
	CardholderAddress     *string    `gorm:"type:varchar(100)" json:"carholder_address"`
	CardholderNationalID  *string    `gorm:"type:varchar(20)" json:"carholder_national_id"`
	CardInfoBinCode       *string    `gorm:"type:varchar(6)" json:"card_info_bin_code"`
	CardInfoLastFour      *string    `gorm:"type:varchar(4)" json:"card_info_last_four"`
	CardInfoIssuer        *string    `gorm:"type:varchar(50)" json:"card_info_issuer"`
	CardInfoFunding       *uint      `gorm:"type:tinyint" json:"card_info_funding"`
	CardInfoType          *uint      `gorm:"type:tinyint" json:"card_info_type"`
	CardInfoLevel         *string    `gorm:"type:varchar(10)" json:"card_info_level"`
	CardInfoCountry       *string    `gorm:"type:varchar(30)" json:"card_info_country"`
	CardInfoCountryCode   *string    `gorm:"type:varchar(10)" json:"card_info_country_code"`
	CardInfoExpiryDate    *string    `gorm:"type:varchar(6);" json:"card_info_expiry_date"`
}

type CardInfo struct {
	BinCode     string `json:"bin_code"`
	LastFour    string `json:"last_four"`
	Issuer      string `json:"issuer"`
	Funding     uint   `json:"funding"`
	Type        uint   `json:"type"`
	Level       string `json:"level"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	ExpiryDate  string `json:"expiry_date"`
}

type Cardholder struct {
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	ZipCode     string `json:"zip_code"`
	Address     string `json:"address"`
	NationalID  string `json:"national_id"`
}

type DonationRecord struct {
	IsPeriodic  bool       `json:"is_periodic"`
	CardInfo    CardInfo   `json:"card_info"`
	Cardholder  Cardholder `json:"cardholder"`
	Amount      uint       `json:"amount"`
	Currency    string     `json:"currency"`
	Details     string     `json:"details"`
	OrderNumber string     `json:"order_number"`
}
