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
	PlainCardInfo
	PlainCardholder
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
	Currency    string     `gorm:"type:char(3);default:'TWD';not null" json:"currency"`
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
	UserID      uint       `gorm:"type:int(10) unsigned;not null" json:"user_id"`
	ZipCode     string     `gorm:"type:varchar(10)" json:"zip_code"`
}

type PeriodicDonation struct {
	Amount        uint        `gorm:"type:int(10) unsigned;not null;index:idx_periodic_donations_amount" json:"amount"`
	CardKey       string      `gorm:"type:tinyblob" json:"card_key"`
	CardToken     string      `gorm:"type:tinyblob" json:"card_token"`
	CreatedAt     time.Time   `json:"created_at"`
	Currency      string      `gorm:"type:varchar(3);default:'TWD';not null" json:"currency"`
	DeletedAt     *time.Time  `json:"deleted_at"`
	ID            uint        `gorm:"primary_key" json:"id"`
	LastSuccessAt null.Time   `json:"last_success_at"`
	SendReceipt   null.String `gorm:"type:ENUM('no', 'monthly', 'yearly');" json:"send_receipt"`
	Status        string      `gorm:"type:ENUM('to_pay','paying','paid','fail');not null" json:"status"`
	UpdatedAt     time.Time   `json:"updated_at"`
	UserID        uint        `gorm:"type:int(10) unsigned;not null" json:"user_id"`
	ToFeedback    null.Bool   `gorm:"type:tinyint(1);default:1" json:"to_feedback"`
	PlainCardInfo
	PlainCardholder
}

// PlainCardInfo is used to store/retrieve data to/from persistent donation storage
type PlainCardInfo struct {
	CardInfoBinCode     null.String `gorm:"type:varchar(6)" json:"card_info_bin_code"`
	CardInfoCountry     null.String `gorm:"type:varchar(30)" json:"card_info_country"`
	CardInfoCountryCode null.String `gorm:"type:varchar(10)" json:"card_info_country_code"`
	CardInfoExpiryDate  null.String `gorm:"type:varchar(6);" json:"card_info_expiry_date"`
	CardInfoFunding     null.Int    `gorm:"type:tinyint" json:"card_info_funding"`
	CardInfoIssuer      null.String `gorm:"type:varchar(50)" json:"card_info_issuer"`
	CardInfoLastFour    null.String `gorm:"type:varchar(4)" json:"card_info_last_four"`
	CardInfoLevel       null.String `gorm:"type:varchar(10)" json:"card_info_level"`
	CardInfoType        null.Int    `gorm:"type:tinyint" json:"card_info_type"`
}

// BinCode is for `copier` to copy from `BinCode`to `CardInfoBinCode`
func (p *PlainCardInfo) BinCode(binCode null.String) {
	p.CardInfoBinCode = binCode
}

// Country is for `copier` to copy from `Country`to `CardInfoCountry`
func (p *PlainCardInfo) Country(country null.String) {
	p.CardInfoCountry = country
}

// CountryCode is for `copier` to copy from `CountryCode`to `CardInfoCountryCode`
func (p *PlainCardInfo) CountryCode(countryCode null.String) {
	p.CardInfoCountryCode = countryCode
}

// ExpiryDate is for `copier` to copy from `ExpiryDate`to `CardInfoExpiryDate`
func (p *PlainCardInfo) ExpiryDate(expiryDate null.String) {
	p.CardInfoExpiryDate = expiryDate
}

// Funding is for `copier` to copy from `Funding`to `CardInfoFunding`
func (p *PlainCardInfo) Funding(funding null.Int) {
	p.CardInfoFunding = funding
}

// Issuer is for `copier` to copy from `Issuer`to `CardInfoIssuer`
func (p *PlainCardInfo) Issuer(issuer null.String) {
	p.CardInfoIssuer = issuer
}

// LastFour is for `copier` to copy from `LastFour`to `CardInfoLastFour`
func (p *PlainCardInfo) LastFour(lastFour null.String) {
	p.CardInfoLastFour = lastFour
}

// Level is for `copier` to copy from `Level`to `CardInfoLevel`
func (p *PlainCardInfo) Level(level null.String) {
	p.CardInfoLevel = level
}

// Type is for `copier` to copy from `Type`to `CardInfoType`
func (p *PlainCardInfo) Type(t null.Int) {
	p.CardInfoType = t
}

// CardInfo is used to communicate with clients
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

// CardInfoBinCode is for `copier` to copy from `CardInfoBinCode`to `BinCode`
func (ci *CardInfo) CardInfoBinCode(binCode null.String) {
	ci.BinCode = binCode
}

// CardInfoCountry is for `copier` to copy from `CardInfoCountry`to `Country`
func (ci *CardInfo) CardInfoCountry(country null.String) {
	ci.Country = country
}

// CardInfoCountryCode is for `copier` to copy from `CardInfoCountryCode`to `CountryCode`
func (ci *CardInfo) CardInfoCountryCode(countryCode null.String) {
	ci.CountryCode = countryCode
}

// CardInfoExpiryDate is for `copier` to copy from `CardInfoExpiryDate`to `ExpiryDate`
func (ci *CardInfo) CardInfoExpiryDate(expiryDate null.String) {
	ci.ExpiryDate = expiryDate
}

// CardInfoFunding is for `copier` to copy from `CardInfoFunding`to `Funding`
func (ci *CardInfo) CardInfoFunding(funding null.Int) {
	ci.Funding = funding
}

// CardInfoIssuer is for `copier` to copy from `CardInfoIssuer`to `Issuer`
func (ci *CardInfo) CardInfoIssuer(issuer null.String) {
	ci.Issuer = issuer
}

// CardInfoLastFour is for `copier` to copy from `CardInfoLastFour`to `LastFour`
func (ci *CardInfo) CardInfoLastFour(lastFour null.String) {
	ci.LastFour = lastFour
}

// CardInfoLevel is for `copier` to copy from `CardInfoLevel`to `Level`
func (ci *CardInfo) CardInfoLevel(level null.String) {
	ci.Level = level
}

// CardInfoType is for `copier` to copy from `CardInfoType`to `Type`
func (ci *CardInfo) CardInfoType(t null.Int) {
	ci.Type = t
}

// PlainCardholder is used to store/retrieve data to/from persistent donation storage
type PlainCardholder struct {
	CardholderAddress     null.String `gorm:"type:varchar(100)" json:"cardholder_address"`
	CardholderEmail       string      `gorm:"type:varchar(100);not null" json:"cardholder_email"`
	CardholderName        null.String `gorm:"type:varchar(30)" json:"cardholder_name"`
	CardholderNationalID  null.String `gorm:"type:varchar(20)" json:"cardholder_national_id"`
	CardholderPhoneNumber null.String `gorm:"type:varchar(20)" json:"cardholder_phone_number"`
	CardholderZipCode     null.String `gorm:"type:varchar(10)" json:"cardholder_zip_code"`
}

// Name is for `copier` to copy from `Name`to `CardholderName`
func (p *PlainCardholder) Name(name null.String) {
	p.CardholderName = name
}

// Email is for `copier` to copy from `Email`to `CardholderEmail`
func (p *PlainCardholder) Email(email string) {
	p.CardholderEmail = email
}

// NationalID is for `copier` to copy from `NationalID` to `CardholderNationalID`
func (p *PlainCardholder) NationalID(nationalID null.String) {
	p.CardholderNationalID = nationalID
}

// Address is for `copier` to copy from `Address` to `CardholderAddress`
func (p *PlainCardholder) Address(addr null.String) {
	p.CardholderAddress = addr
}

// PhoneNumber is for `copier` to copy from `PhoneNumber` to `CardholderPhoneNumber`
func (p *PlainCardholder) PhoneNumber(phoneNumber null.String) {
	p.CardholderPhoneNumber = phoneNumber
}

// ZipCode is for `copier` to copy from `ZipCode` to `CardholderZipCode`
func (p *PlainCardholder) ZipCode(zipCode null.String) {
	p.CardholderZipCode = zipCode
}

// Cardholder is used to communicate with clients
type Cardholder struct {
	Address     null.String `json:"address"`
	Email       string      `json:"email" binding:"required,email"`
	Name        null.String `json:"name"`
	NationalID  null.String `json:"national_id"`
	PhoneNumber null.String `json:"phone_number"`
	ZipCode     null.String `json:"zip_code"`
}

// CardholderAddress is for `copier` to copy from `CardholderAddress` to `Address`
func (c *Cardholder) CardholderAddress(addr null.String) {
	c.Address = addr
}

// CardholderEmail is for `copier` to copy from `CardholderEmail` to `Email`
func (c *Cardholder) CardholderEmail(email string) {
	c.Email = email
}

// CardholderName is for `copier` to copy from `CardholderName` to `Name`
func (c *Cardholder) CardholderName(name null.String) {
	c.Name = name
}

// CardholderNationalID is for `copier` to copy from `CardholderNationalID` to `NationalID`
func (c *Cardholder) CardholderNationalID(nationalID null.String) {
	c.NationalID = nationalID
}

// CardholderPhoneNumber is for `copier` to copy from `CardholderPhoneNumber` to `PhoneNumber`
func (c *Cardholder) CardholderPhoneNumber(phoneNumber null.String) {
	c.PhoneNumber = phoneNumber
}

// CardholderZipCode is for `copier` to copy from `CardholderZipCode` to `ZipCode`
func (c *Cardholder) CardholderZipCode(zipCode null.String) {
	c.ZipCode = zipCode
}
