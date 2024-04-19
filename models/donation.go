package models

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

type TappayResp struct {
	Acquirer                 string      `gorm:"type:varchar(50);not null" json:"acquirer"`
	AuthCode                 string      `gorm:"type:varchar(6);not null" json:"auth_code"`
	BankResultCode           null.String `gorm:"type:varchar(50)" json:"bank_result_code"`
	BankResultMsg            null.String `gorm:"type:varchar(128)" json:"bank_result_msg"`
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
	Email              string      `gorm:"column:cardholder_email;type:varchar(100);not null" json:"email" binding:"omitempty,email"`
	Name               null.String `gorm:"column:cardholder_name;type:varchar(30)" json:"name"`
	FirstName          null.String `gorm:"column:cardholder_first_name;type:varchar(30)" json:"first_name"`
	LastName           null.String `gorm:"column:cardholder_last_name;type:varchar(30)" json:"last_name"`
	Nickname           null.String `gorm:"column:cardholder_nickname;type:varchar(50)" json:"nickname"`
	Title              null.String `gorm:"column:cardholder_title;type:varchar(30)" json:"title"`
	LegalName          null.String `gorm:"column:cardholder_legal_name;type:varchar(50)" json:"legal_name"`
	Gender             null.String `gorm:"column:cardholder_gender;type:varchar(2)" json:"gender"` // e.g., "M", "F", "X, "U"
	AgeRange           null.String `gorm:"column:cardholder_age_range;type:ENUM('less_than_18', '18_to_24', '25_to_34', '35_to_44', '45_to_54', '55_to_64', 'above_65')" json:"age_range"`
	ReadPreference     null.String `gorm:"column:cardholder_read_preference;type:SET('international', 'cross_straits', 'human_right', 'society', 'environment', 'education', 'politics', 'economy', 'culture', 'art', 'life', 'health', 'sport', 'all')" json:"read_preference"` // e.g. "international, art, sport"
	WordsForTwreporter null.String `gorm:"column:cardholder_words_for_twreporter;type:varchar(255)" json:"words_for_twreporter"`
	NationalID         null.String `gorm:"column:cardholder_national_id;type:varchar(20)" json:"national_id"`
	SecurityID         null.String `gorm:"column:cardholder_security_id;type:varchar(20)" json:"security_id"`
	PhoneNumber        null.String `gorm:"column:cardholder_phone_number;type:varchar(20)" json:"phone_number"`
	ZipCode            null.String `gorm:"column:cardholder_zip_code;type:varchar(10)" json:"zip_code"`
	Address            null.String `gorm:"column:cardholder_address;type:varchar(100)" json:"address"`
	AddressCountry     null.String `gorm:"column:cardholder_address_country;type:varchar(45)" json:"address_country"`
	AddressState       null.String `gorm:"column:cardholder_address_state;type:varchar(45)" json:"address_state"`
	AddressCity        null.String `gorm:"column:cardholder_address_city;type:varchar(45)" json:"address_city"`
	AddressDetail      null.String `gorm:"column:cardholder_address_detail;type:varchar(255)" json:"address_detail"`
	AddressZipCode     null.String `gorm:"column:cardholder_address_zip_code;type:varchar(10)" json:"address_zip_code"`
}

type Receipt struct {
	Header         null.String `gorm:"column:receipt_header;type:varchar(128)" json:"header"`
	SecurityID     null.String `gorm:"column:receipt_security_id;type:varchar(20)" json:"security_id"`
	Email          null.String `gorm:"column:receipt_email;type:varchar(100)" json:"email"`
	AddressCountry null.String `gorm:"column:receipt_address_country;type:varchar(45)" json:"address_country"`
	AddressState   null.String `gorm:"column:receipt_address_state;type:varchar(45)" json:"address_state"`
	AddressCity    null.String `gorm:"column:receipt_address_city;type:varchar(45)" json:"address_city"`
	AddressDetail  null.String `gorm:"column:receipt_address_detail;type:varchar(255)" json:"address_detail"`
	AddressZipCode null.String `gorm:"column:receipt_address_zip_code;type:varchar(10)" json:"address_zip_code"`
}

// https://docs.tappaysdk.com/tutorial/zh/back.html#request-body pay_info
// masked_credit_card_number will be preprocessed and stored in the CardInfo.LastFour
type PayInfo struct {
	Method                 null.String `gorm:"column:linepay_method;type:ENUM('CREDIT_CARD', 'BALANCE', 'POINT')" json:"method"`
	MaskedCreditCardNumber null.String `gorm:"-" json:"masked_credit_card_number"`
	Point                  null.Int    `gorm:"column:linepay_point;type:int;" json:"point"`
}

type PayByPrimeDonation struct {
	CardInfo
	Cardholder
	TappayResp
	Receipt
	PayInfo          `json:"pay_info"`
	Amount           uint       `gorm:"not null" json:"amount"`
	CreatedAt        time.Time  `json:"created_at"`
	Currency         string     `gorm:"type:varchar(3);default:'TWD';not null" json:"currency"`
	DeletedAt        *time.Time `json:"deleted_at"`
	Details          string     `gorm:"type:varchar(50);not null" json:"details"`
	ID               uint       `gorm:"primary_key" json:"id"`
	MerchantID       string     `gorm:"type:varchar(30);not null" json:"merchant_id"`
	Notes            string     `gorm:"type:varchar(100)" json:"notes"`
	OrderNumber      string     `gorm:"type:varchar(50);not null" json:"order_number"`
	PayMethod        string     `gorm:"type:ENUM('credit_card','line','apple','google','samsung');not null;index:idx_pay_by_prime_donations_cardholder_email_pay_method" json:"pay_method"`
	SendReceipt      string     `gorm:"type:ENUM('yearly', 'monthly', 'no', 'no_receipt', 'digital_receipt_by_month', 'digital_receipt_by_year', 'paperback_receipt_by_month', 'paperback_receipt_by_year');default:'no_receipt'" json:"send_receipt"`
	Status           string     `gorm:"type:ENUM('paying','paid','fail','refunded');not null" json:"status"`
	UpdatedAt        time.Time  `json:"updated_at"`
	UserID           uint       `gorm:"type:int(10);unsigned;not null" json:"user_id"`
	IsAnonymous      null.Bool  `gorm:"type:tinyint(1);default:0" json:"is_anonymous"`
	AutoTaxDeduction null.Bool  `gorm:"type:tinyint(1)" json:"auto_tax_deduction"`
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
	Status      string     `gorm:"type:ENUM('paying','paid','fail','refunded');not null" json:"status"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type PayByOtherMethodDonation struct {
	Address       string      `gorm:"type:varchar(100)" json:"address"`
	Amount        uint        `gorm:"type:int(10) unsigned;index:idx_pay_by_other_donations_amount" json:"amount"`
	CreatedAt     time.Time   `json:"created_at"`
	Currency      string      `gorm:"type:char(3);default:'TWD';not null" json:"currency"`
	DeletedAt     *time.Time  `json:"deleted_at"`
	Details       string      `gorm:"type:varchar(50);not null" json:"details"`
	Email         string      `gorm:"type:varchar(100);not null" json:"email"`
	ID            uint        `gorm:"primary_key" json:"id"`
	MerchantID    string      `gorm:"type:varchar(30);not null" json:"merchant_id"`
	Name          string      `gorm:"type:varchar(30)" json:"name"`
	NationalID    string      `gorm:"type:varchar(20)" json:"national_id"`
	SecurityID    string      `gorm:"type:varchar(20)" json:"security_id"`
	Notes         string      `gorm:"type:varchar(100)" json:"notes"`
	OrderNumber   string      `gorm:"type:varchar(50);not null" json:"order_number"`
	PayMethod     string      `gorm:"type:varchar(50);not null;index:idx_pay_by_other_donations_pay_method" json:"pay_method"`
	PhoneNumber   string      `gorm:"type:varchar(20)" json:"phone_number"`
	SendReceipt   string      `gorm:"type:ENUM('no', 'yearly', 'monthly', 'no_receipt', 'digital_receipt_by_month', 'digital_receipt_by_year', 'paperback_receipt_by_month', 'paperback_receipt_by_year');default:'no_receipt'" json:"send_receipt"`
	UpdatedAt     time.Time   `json:"updated_at"`
	UserID        uint        `gorm:"type:int(10) unsigned;not null" json:"user_id"`
	ZipCode       string      `gorm:"type:varchar(10)" json:"zip_code"`
	ReceiptHeader null.String `gorm:"type:varchar(128)" json:"receipt_header"`
}

type PeriodicDonation struct {
	Cardholder
	CardInfo
	Receipt
	Amount           uint       `gorm:"type:int(10) unsigned;not null;index:idx_periodic_donations_amount" json:"amount"`
	CardKey          string     `gorm:"type:tinyblob" json:"card_key"`
	CardToken        string     `gorm:"type:tinyblob" json:"card_token"`
	CreatedAt        time.Time  `json:"created_at"`
	Currency         string     `gorm:"type:varchar(3);default:'TWD';not null" json:"currency"`
	DeletedAt        *time.Time `json:"deleted_at"`
	Details          string     `gorm:"type:varchar(50);not null" json:"details"`
	Frequency        string     `gorm:"type:ENUM('monthly', 'yearly');default:'monthly'" json:"frequency"`
	ID               uint       `gorm:"primary_key" json:"id"`
	LastSuccessAt    null.Time  `json:"last_success_at"`
	MaxPaidTimes     uint       `json:"max_paid_times" gorm:"type:int;not null;default:2147483647"`
	Notes            string     `gorm:"type:varchar(100)" json:"notes"`
	OrderNumber      string     `gorm:"type:varchar(50);not null" json:"order_number"`
	SendReceipt      string     `gorm:"type:ENUM('yearly', 'no', 'no_receipt', 'digital_receipt_by_year', 'paperback_receipt_by_year');default:'no_receipt'" json:"send_receipt"`
	Status           string     `gorm:"type:ENUM('to_pay','paying','paid','fail','stopped','invalid');not null" json:"status"`
	ToFeedback       null.Bool  `gorm:"type:tinyint(1);default:1" json:"to_feedback"`
	UpdatedAt        time.Time  `json:"updated_at"`
	UserID           uint       `gorm:"type:int(10) unsigned;not null" json:"user_id"`
	IsAnonymous      null.Bool  `gorm:"type:tinyint(1);default:0" json:"is_anonymous"`
	AutoTaxDeduction null.Bool  `gorm:"type:tinyint(1)" json:"auto_tax_deduction"`
	PayMethod        string     `gorm:"type:ENUM('credit_card','line','apple','google','samsung')" json:"pay_method"`
}

type GeneralDonation struct {
	ID               uint        `json:"id"`
	Type             string      `json:"type"`
	Amount           uint        `json:"amount"`
	CreatedAt        time.Time   `json:"created_at"`
	OrderNumber      string      `json:"order_number"`
	SendReceipt      string      `json:"send_receipt"`
	Status           string      `json:"status"`
	PayMethod        string      `json:"pay_method"`
	BinCode          null.String `gorm:"column:card_info_bin_code" json:"bin_code,omitempty"`
	CardLastFour     null.String `gorm:"column:card_info_last_four" json:"card_last_four, omitempty"`
	CardType         null.String `gorm:"column:card_info_type" json:"card_type, omitempty"`
	IsAnonymous      null.Bool  `gorm:"type:tinyint(1);default:0" json:"is_anonymous"`
	FirstName        null.String `gorm:"column:cardholder_first_name" json:"first_name,omitempty"`
	LastName         null.String `gorm:"column:cardholder_last_name" json:"last_name,omitempty"`
	Header           null.String `gorm:"column:receipt_header" json:"receipt_header,omitempty"`
	AddressCountry   null.String `gorm:"column:receipt_address_country" json:"address_country,omitempty"`
	AddressState     null.String `gorm:"column:receipt_address_state" json:"address_state,omitempty"`
	AddressCity      null.String `gorm:"column:receipt_address_city" json:"address_city,omitempty"`
	AddressDetail    null.String `gorm:"column:receipt_address_detail" json:"address_detail,omitempty"`
	AddressZipCode   null.String `gorm:"column:receipt_address_zip_code" json:"address_zip_code,omitempty"`
}

type Payment struct {
	CreatedAt        time.Time   `json:"created_at"`
	OrderNumber      string      `json:"order_number"`
	Status           string      `json:"status"`
	Amount           uint        `json:"amount"`
}
