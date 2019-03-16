package models

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

type TappayResp struct {
	Acquirer                 string      `gorm:"type:varchar(50);not null" json:"acquirer"`
	AuthCode                 string      `gorm:"type:varchar(6);not null" json:"auth_code"`
	BankResultCode           null.String `gorm:"type:varchar(50)" json:"bank_result_code"`
	BankResultMsg            null.String `gorm:"type:varchar(50)" json:"bank_result_msg"`
	BankTransactionEndTime   null.Time   `json:"bank_transaction_end_time"`
	BankTransactionID        string      `gorm:"type:varchar(50);not null" json:"bank_transaction_id"`
	BankTransactionStartTime null.Time   `json:"bank_transaction_start_time"`
	Msg                      string      `gorm:"type:varchar(100);not null" json:"msg"`
	RecTradeID               string      `gorm:"type:varchar(20);not null" json:"rec_trade_id"`
	TappayApiStatus          null.Int    `json:"tappay_api_status"`
	TappayRecordStatus       null.Int    `json:"tappay_record_status"`
	TransactionTime          null.Time   `json:"transaction_time"`
}

type CardInfo struct {
	BinCode     null.String `gorm:"column:card_info_bin_code;type:varchar(6)" json:"bin_code"`
	Country     null.String `gorm:"column:card_info_country;type:varchar(30)" json:"country"`
	CountryCode null.String `gorm:"column:card_info_country_code;type:varchar(10)" json:"country_code"`
	ExpiryDate  null.String `gorm:"column:card_info_expiry_date;type:varchar(6);" json:"expiry_date"`
	Funding     null.Int    `gorm:"column:card_info_funding;type:tinyint" json:"funding"`
	Issuer      null.String `gorm:"column:card_info_issuer;type:varchar(50)" json:"issuer"`
	LastFour    null.String `gorm:"column:card_info_last_four;type:varchar(4)" json:"last_four"`
	Level       null.String `gorm:"column:card_info_level;type:varchar(10)" json:"level"`
	Type        null.Int    `gorm:"column:card_info_type;type:tinyint" json:"type"`
}

type Cardholder struct {
	Address     null.String `gorm:"column:cardholder_address;type:varchar(100)" json:"address"`
	Email       string      `gorm:"column:cardholder_email;type:varchar(100);not null" json:"email" binding:"omitempty,email"`
	Name        null.String `gorm:"column:cardholder_name;type:varchar(30)" json:"name"`
	NationalID  null.String `gorm:"column:cardholder_national_id;type:varchar(20)" json:"national_id"`
	PhoneNumber null.String `gorm:"column:cardholder_phone_number;type:varchar(20)" json:"phone_number"`
	ZipCode     null.String `gorm:"column:cardholder_zip_code;type:varchar(10)" json:"zip_code"`
}

type PayByPrimeDonation struct {
	CardInfo
	Cardholder
	TappayResp
	Amount      uint       `gorm:"not null" json:"amount"`
	CreatedAt   time.Time  `json:"created_at"`
	Currency    string     `gorm:"type:varchar(3);default:'TWD';not null" json:"currency"`
	DeletedAt   *time.Time `json:"deleted_at"`
	Details     string     `gorm:"type:varchar(50);not null" json:"details"`
	ID          uint       `gorm:"primary_key" json:"id"`
	MerchantID  string     `gorm:"type:varchar(30);not null" json:"merchant_id"`
	Notes       string     `gorm:"type:varchar(100)" json:"notes"`
	OrderNumber string     `gorm:"type:varchar(50);not null" json:"order_number"`
	PayMethod   string     `gorm:"type:ENUM('credit_card','line','apple','google','samsung');not null;index:idx_pay_by_prime_donations_cardholder_email_pay_method" json:"pay_method"`
	SendReceipt string     `gorm:"type:ENUM('no', 'monthly');default:'monthly'" json:"send_receipt"`
	Status      string     `gorm:"type:ENUM('paying','paid','fail');not null" json:"status"`
	UpdatedAt   time.Time  `json:"updated_at"`
	UserID      uint       `gorm:"type:int(10);unsigned;not null" json:"user_id"`
}

type PayByCardTokenDonation struct {
	TappayResp
	Amount      uint       `gorm:"not null;index:idx_pay_by_card_token_donations_amount" json:"amount"`
	CreatedAt   time.Time  `json:"created_at"`
	Currency    string     `gorm:"type:varchar(3);default:'TWD';not null" json:"currency"`
	DeletedAt   *time.Time `json:"deleted_at"`
	Details     string     `gorm:"type:varchar(50);not null" json:"details"`
	ID          uint       `gorm:"primary_key" json:"id"`
	MerchantID  string     `gorm:"type:varchar(30);not null" json:"merchant_id"`
	OrderNumber string     `gorm:"type:varchar(50);not null" json:"order_number"`
	PeriodicID  uint       `gorm:"not null;index:idx_pay_by_card_token_donations_periodic_id" json:"periodic_id"`
	Status      string     `gorm:"type:ENUM('paying','paid','fail');not null" json:"status"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type PayByOtherMethodDonation struct {
	Address     string     `gorm:"type:varchar(100)" json:"address"`
	Amount      uint       `gorm:"type:int(10) unsigned;index:idx_pay_by_other_donations_amount" json:"amount"`
	CreatedAt   time.Time  `json:"created_at"`
	Currency    string     `gorm:"type:char(3);default:'TWD';not null" json:"currency"`
	DeletedAt   *time.Time `json:"deleted_at"`
	Details     string     `gorm:"type:varchar(50);not null" json:"details"`
	Email       string     `gorm:"type:varchar(100);not null" json:"email"`
	ID          uint       `gorm:"primary_key" json:"id"`
	MerchantID  string     `gorm:"type:varchar(30);not null" json:"merchant_id"`
	Name        string     `gorm:"type:varchar(30)" json:"name"`
	NationalID  string     `gorm:"type:varchar(20)" json:"national_id"`
	Notes       string     `gorm:"type:varchar(100)" json:"notes"`
	OrderNumber string     `gorm:"type:varchar(50);not null" json:"order_number"`
	PayMethod   string     `gorm:"type:varchar(50);not null;index:idx_pay_by_other_donations_pay_method" json:"pay_method"`
	PhoneNumber string     `gorm:"type:varchar(20)" json:"phone_number"`
	SendReceipt string     `gorm:"type:ENUM('no', 'monthly');default:'monthly'" json:"send_receipt"`
	UpdatedAt   time.Time  `json:"updated_at"`
	UserID      uint       `gorm:"type:int(10) unsigned;not null" json:"user_id"`
	ZipCode     string     `gorm:"type:varchar(10)" json:"zip_code"`
}

type PeriodicDonation struct {
	Cardholder
	CardInfo
	Amount        uint       `gorm:"type:int(10) unsigned;not null;index:idx_periodic_donations_amount" json:"amount"`
	CardKey       string     `gorm:"type:tinyblob" json:"card_key"`
	CardToken     string     `gorm:"type:tinyblob" json:"card_token"`
	CreatedAt     time.Time  `json:"created_at"`
	Currency      string     `gorm:"type:varchar(3);default:'TWD';not null" json:"currency"`
	DeletedAt     *time.Time `json:"deleted_at"`
	Details       string     `gorm:"type:varchar(50);not null" json:"details"`
	Frequency     string     `gorm:"type:ENUM('monthly', 'yearly');default:'monthly'" json:"frequency"`
	ID            uint       `gorm:"primary_key" json:"id"`
	LastSuccessAt null.Time  `json:"last_success_at"`
	MaxPaidTimes  uint       `json:"max_paid_times" gorm:"type:int;not null;default:2147483647"`
	Notes         string     `gorm:"type:varchar(100)" json:"notes"`
	OrderNumber   string     `gorm:"type:varchar(50);not null" json:"order_number"`
	SendReceipt   string     `gorm:"type:ENUM('no', 'monthly', 'yearly');default:'monthly'" json:"send_receipt"`
	Status        string     `gorm:"type:ENUM('to_pay','paying','paid','fail','stopped','invalid');not null" json:"status"`
	ToFeedback    null.Bool  `gorm:"type:tinyint(1);default:1" json:"to_feedback"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UserID        uint       `gorm:"type:int(10) unsigned;not null" json:"user_id"`
}
