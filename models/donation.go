package models

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

type PayByPrimeDonation struct {
	Acquirer                 string      `gorm:"type:varchar(50);not null" json:"acquirer"`
	Amount                   uint        `gorm:"not null" json:"amount"`
	AuthCode                 string      `gorm:"type:varchar(6);not null" json:"auth_code"`
	BankResultCode           null.String `gorm:"type:varchar(50)" json:"bank_result_code"`
	BankResultMsg            null.String `gorm:"type:varchar(50)" json:"bank_result_msg"`
	BankTransactionEndTime   null.Time   `json:"bank_transaction_end_time"`
	BankTransactionID        string      `gorm:"type:varchar(50);not null" json:"bank_transaction_id"`
	BankTransactionStartTime null.Time   `json:"bank_transaction_start_time"`
	CardInfoBinCode          null.String `gorm:"type:varchar(6)" json:"card_info_bin_code"`
	CardInfoCountry          null.String `gorm:"type:varchar(30)" json:"card_info_country"`
	CardInfoCountryCode      null.String `gorm:"type:varchar(10)" json:"card_info_country_code"`
	CardInfoExpiryDate       null.String `gorm:"type:varchar(6);" json:"card_info_expiry_date"`
	CardInfoFunding          null.Int    `gorm:"type:tinyint" json:"card_info_funding"`
	CardInfoIssuer           null.String `gorm:"type:varchar(50)" json:"card_info_issuer"`
	CardInfoLastFour         null.String `gorm:"type:varchar(4)" json:"card_info_last_four"`
	CardInfoLevel            null.String `gorm:"type:varchar(10)" json:"card_info_level"`
	CardInfoType             null.Int    `gorm:"type:tinyint" json:"card_info_type"`
	CardholderAddress        null.String `gorm:"type:varchar(100)" json:"cardholder_address"`
	CardholderEmail          string      `gorm:"type:varchar(100);not null;index:idx_pay_by_prime_donations_cardholder_email_pay_method" json:"cardholder_email"`
	CardholderName           null.String `gorm:"type:varchar(30)" json:"cardholder_name"`
	CardholderNationalID     null.String `gorm:"type:varchar(20)" json:"cardholder_national_id"`
	CardholderPhoneNumber    null.String `gorm:"type:varchar(20)" json:"cardholder_phone_number"`
	CardholderZipCode        null.String `gorm:"type:varchar(10)" json:"cardholder_zip_code"`
	CreatedAt                time.Time   `json:"created_at"`
	Currency                 string      `gorm:"type:varchar(3);default:'TWD';not null" json:"currency"`
	DeletedAt                *time.Time  `json:"deleted_at"`
	Details                  string      `gorm:"type:varchar(50);not null" json:"details"`
	ID                       uint        `gorm:"primary_key" json:"id"`
	MerchantID               string      `gorm:"type:varchar(30);not null" json:"merchant_id"`
	Msg                      string      `gorm:"type:varchar(100);not null" json:"msg"`
	OrderNumber              string      `gorm:"type:varchar(50);not null" json:"order_number"`
	PayMethod                string      `gorm:"type:ENUM('credit_card','line','apple','google','samsung');not null;index:idx_pay_by_prime_donations_cardholder_email_pay_method" json:"pay_method"`
	RecTradeID               string      `gorm:"type:varchar(20);not null" json:"rec_trade_id"`
	SendReceipt              null.String `gorm:"type:ENUM('no', 'monthly');" json:"send_receipt"`
	Status                   string      `gorm:"type:ENUM('paying','paid','fail');not null" json:"status"`
	TappayApiStatus          null.Int    `json:"tappay_api_status"`
	TappayRecordStatus       null.Int    `json:"tappay_record_status"`
	TransactionTime          null.Time   `json:"transaction_time"`
	UpdatedAt                time.Time   `json:"updated_at"`
	UserID                   uint        `gorm:"type:int(10);unsigned;not null" json:"user_id"`
}

type PayByCardTokenDonation struct {
	Acquirer                 string      `gorm:"type:varchar(50);not null" json:"acquirer"`
	Amount                   uint        `gorm:"not null;index:idx_pay_by_card_token_donations_amount" json:"amount"`
	AuthCode                 string      `gorm:"type:varchar(6);not null" json:"auth_code"`
	BankResultCode           null.String `gorm:"type:varchar(50)" json:"bank_result_code"`
	BankResultMsg            null.String `gorm:"type:varchar(50)" json:"bank_result_msg"`
	BankTransactionEndTime   null.Time   `json:"bank_transaction_end_time"`
	BankTransactionID        string      `gorm:"type:varchar(50);not null" json:"bank_transaction_id"`
	BankTransactionStartTime null.Time   `json:"bank_transaction_start_time"`
	CreatedAt                time.Time   `json:"created_at"`
	Currency                 string      `gorm:"type:varchar(3);default:'TWD';not null" json:"currency"`
	DeletedAt                *time.Time  `json:"deleted_at"`
	Details                  string      `gorm:"type:varchar(50);not null" json:"details"`
	ID                       uint        `gorm:"primary_key" json:"id"`
	MerchantID               string      `gorm:"type:varchar(30);not null" json:"merchant_id"`
	Msg                      string      `gorm:"type:varchar(100);not null" json:"msg"`
	OrderNumber              string      `gorm:"type:varchar(50);not null" json:"order_number"`
	PeriodicID               uint        `gorm:"not null;index:idx_pay_by_card_token_donations_periodic_id" json:"periodic_id"`
	RecTradeID               string      `gorm:"type:varchar(20);not null" json:"rec_trade_id"`
	Status                   string      `gorm:"type:ENUM('paying','paid','fail');not null" json:"status"`
	TappayApiStatus          null.Int    `json:"tappay_api_status"`
	TappayRecordStatus       null.Int    `json:"tappay_record_status"`
	TransactionTime          null.Time   `json:"transaction_time"`
	UpdatedAt                time.Time   `json:"updated_at"`
}

type PayByOtherMethodDonation struct {
	Address     string     `gorm:"type:varchar(100)" json:"address"`
	Amount      uint       `gorm:"type:int(10) unsigned;index:idx_pay_by_other_donations_amount" json:"amount"`
	CreatedAt   time.Time  `json:"created_at"`
	Curreny     string     `gorm:"type:char(3);default:'TWD';not null" json:"currency"`
	DeletedAt   *time.Time `json:"deleted_at"`
	Details     string     `gorm:"type:varchar(50);not null" json:"details"`
	Email       string     `gorm:"type:varchar(100);not null" json:"email"`
	ID          uint       `gorm:"primary_key" json:"id"`
	MerchantID  string     `gorm:"type:varchar(30);not null" json:"merchant_id"`
	Name        string     `gorm:"type:varchar(30)" json:"name"`
	NationalID  string     `gorm:"type:varchar(20)" json:"national_id"`
	OrderNumber string     `gorm:"type:varchar(50);not null" json:"order_number"`
	PayMethod   string     `gorm:"type:varchar(50);not null;index:idx_pay_by_other_donations_pay_method" json:"pay_method"`
	SendReceipt string     `gorm:"type:ENUM('no', 'monthly');" json:"send_receipt"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ZipCode     string     `gorm:"type:varchar(10)" json:"zip_code"`
}

type PeriodicDonation struct {
	Amount                uint        `gorm:"type:int(10) unsigned;not null;index:idx_periodic_donations_amount" json:"amount"`
	CardInfoBinCode       null.String `gorm:"type:varchar(6)" json:"card_info_bin_code"`
	CardInfoCountry       null.String `gorm:"type:varchar(30)" json:"card_info_country"`
	CardInfoCountryCode   null.String `gorm:"type:varchar(10)" json:"card_info_country_code"`
	CardInfoExpiryDate    null.String `gorm:"type:varchar(6);" json:"card_info_expiry_date"`
	CardInfoFunding       null.Int    `gorm:"type:tinyint" json:"card_info_funding"`
	CardInfoIssuer        null.String `gorm:"type:varchar(50)" json:"card_info_issuer"`
	CardInfoLastFour      null.String `gorm:"type:varchar(4)" json:"card_info_last_four"`
	CardInfoLevel         null.String `gorm:"type:varchar(10)" json:"card_info_level"`
	CardInfoType          null.Int    `gorm:"type:tinyint" json:"card_info_type"`
	CardKey               string      `gorm:"type:tinyblob" json:"card_key"`
	CardToken             string      `gorm:"type:tinyblob" json:"card_token"`
	CardholderAddress     null.String `gorm:"type:varchar(100)" json:"cardholder_address"`
	CardholderEmail       string      `gorm:"type:varchar(100);not null" json:"cardholder_email"`
	CardholderName        null.String `gorm:"type:varchar(30)" json:"cardholder_name"`
	CardholderNationalID  null.String `gorm:"type:varchar(20)" json:"cardholder_national_id"`
	CardholderPhoneNumber null.String `gorm:"type:varchar(20)" json:"cardholder_phone_number"`
	CardholderZipCode     null.String `gorm:"type:varchar(10)" json:"cardholder_zip_code"`
	CreatedAt             time.Time   `json:"created_at"`
	Currency              string      `gorm:"type:varchar(3);default:'TWD';not null" json:"currency"`
	DeletedAt             *time.Time  `json:"deleted_at"`
	ID                    uint        `gorm:"primary_key" json:"id"`
	LastSuccessAt         null.Time   `json:"last_success_at"`
	SendReceipt           null.String `gorm:"type:ENUM('no', 'monthly', 'yearly');" json:"send_receipt"`
	Status                string      `gorm:"type:ENUM('to_pay','paying','paid','fail');not null" json:"status"`
	UpdatedAt             time.Time   `json:"updated_at"`
	UserID                uint        `gorm:"type:int(10) unsigned;not null" json:"user_id"`
	ToFeedback            null.Bool   `gorm:"type:tinyint(1);default:1" json:"to_feedback"`
}

type CardInfo struct {
	BinCode     null.String `json:"bin_code"`
	Country     null.String `json:"country"`
	CountryCode null.String `json:"country_code"`
	ExpiryDate  null.String `json:"expiry_date"`
	Funding     null.Int    `json:"funding"`
	Issuer      null.String `json:"issuer"`
	LastFour    null.String `json:"last_four"`
	Level       null.String `json:"level"`
	Type        null.Int    `json:"type"`
}

type Cardholder struct {
	Address     null.String `json:"address"`
	Email       string      `json:"email" binding:"required,email"`
	Name        null.String `json:"name"`
	NationalID  null.String `json:"national_id"`
	PhoneNumber null.String `json:"phone_number"`
	ZipCode     null.String `json:"zip_code"`
}

type DonationRecord struct {
	Amount      uint       `json:"amount"`
	CardInfo    CardInfo   `json:"card_info"`
	Cardholder  Cardholder `json:"cardholder"`
	Currency    string     `json:"currency"`
	Details     string     `json:"details"`
	IsPeriodic  bool       `json:"is_periodic"`
	OrderNumber string     `json:"order_number"`
}
