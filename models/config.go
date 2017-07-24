package models

import (
	"time"
)

const (
	// Office360 global string constants
	Office360 = "office360"

	// ConnSecurityPlain global string constants
	ConnSecurityPlain = "PLAIN"

	// ConnSecurityTLS global string constants
	ConnSecurityTLS = "TLS"

	// ConnSecurityStarttls global string constants
	ConnSecurityStarttls = "STARTTLS"

	// AppSettingsDefaultPath default path of app
	AppSettingsDefaultPath = "http://testtest.twreporter.org:8080"

	// AppSettingsDefaultExpiration default expiration of app
	AppSettingsDefaultExpiration = 168 // 7 days

	// AppSettingsDefaultVersion default version
	AppSettingsDefaultVersion = "v1"

	// EmailSettingsDefaultSMTPServer default smtp server hostname
	EmailSettingsDefaultSMTPServer = "smtp.office365.com"

	// EmailSettingsDefaultSMTPPort default port of smtp
	EmailSettingsDefaultSMTPPort = "587"

	// EmailSettingsDefaultSMTPServerOwner default owner of smtp server
	EmailSettingsDefaultSMTPServerOwner = Office360

	// EmailSettingsDefaultConnSecurity default connection security
	EmailSettingsDefaultConnSecurity = ConnSecurityStarttls

	// DbSettingsDefaultName default database name
	DbSettingsDefaultName = "test_membership"

	// DbSettingsDefaultAddress default address of database
	DbSettingsDefaultAddress = "127.0.0.1"

	// DbSettingsDefaultPort default port of database
	DbSettingsDefaultPort = "3306"

	// EncryptSettingsDefaultSalt default salt for encryption
	EncryptSettingsDefaultSalt = "@#$%"

	// AmazonMailSettingsDefaultCharSet sets the default charset
	AmazonMailSettingsDefaultCharSet = "UTF-8"
)

// AppSettings could be defined in configs/config.json
type AppSettings struct {
	Path       string
	Version    string
	Token      string
	Expiration time.Duration
}

// EmailSettings could be defined in configs/config.json
type EmailSettings struct {
	SMTPUsername       string
	SMTPPassword       string
	SMTPServer         string
	SMTPPort           string
	ConnectionSecurity string
	SMTPServerOwner    string
	FeedbackName       string
	FeedbackEmail      string
}

// AmazonMailSettings could be defined in configs/config.json
type AmazonMailSettings struct {
	Sender    string
	AwsRegion string
	CharSet   string
}

// DBSettings could be defined in configs/config.json
type DBSettings struct {
	Name     string
	User     string
	Password string
	Address  string
	Port     string
}

// MongoDBSettings ...
type MongoDBSettings struct {
	URL     string
	DBName  string
	Timeout int
}

// FacebookSettings could be defined in configs/config.json
type FacebookSettings struct {
	ID       string
	Secret   string
	URL      string
	Statestr string
}

// GoogleSettings could be defined in configs/config.json
type GoogleSettings struct {
	ID       string
	Secret   string
	URL      string
	Statestr string
}

// OauthSettings this contains FacebookSettings and GoogleSettings
type OauthSettings struct {
	FacebookSettings FacebookSettings
	GoogleSettings   GoogleSettings
}

// AlgoliaSettings ...
type AlgoliaSettings struct {
	ApplicationID string
	APIKey        string
}

// ConsumerSettings describes who uses this api
type ConsumerSettings struct {
	Domain   string
	Protocal string
	Host     string
	Port     string
}

// EncryptSettings could be defined in configs/config.json
type EncryptSettings struct {
	Salt string
}

// Config contains all the other configs
type Config struct {
	AlgoliaSettings    AlgoliaSettings
	AppSettings        AppSettings
	EmailSettings      EmailSettings
	AmazonMailSettings AmazonMailSettings
	DBSettings         DBSettings
	MongoDBSettings    MongoDBSettings
	OauthSettings      OauthSettings
	ConsumerSettings   ConsumerSettings
	EncryptSettings    EncryptSettings
}

// SetDefaults could set default value in the Config struct
func (o *Config) SetDefaults() {
	if o.AppSettings.Expiration == 0 {
		o.AppSettings.Expiration = AppSettingsDefaultExpiration
	}
	if o.AppSettings.Version == "" {
		o.AppSettings.Version = AppSettingsDefaultVersion
	}
	if o.DBSettings.Name == "" {
		o.DBSettings.Name = DbSettingsDefaultName
	}
	if o.DBSettings.Address == "" {
		o.DBSettings.Address = DbSettingsDefaultAddress
	}
	if o.DBSettings.Port == "" {
		o.DBSettings.Port = DbSettingsDefaultPort
	}
	if o.EmailSettings.SMTPServer == "" {
		o.EmailSettings.SMTPServer = EmailSettingsDefaultSMTPServer
	}
	if o.EmailSettings.SMTPPort == "" {
		o.EmailSettings.SMTPPort = EmailSettingsDefaultSMTPPort
	}
	if o.EmailSettings.SMTPServerOwner == "" {
		o.EmailSettings.SMTPServerOwner = EmailSettingsDefaultSMTPServerOwner
	}
	if o.EmailSettings.ConnectionSecurity == "" {
		o.EmailSettings.ConnectionSecurity = EmailSettingsDefaultConnSecurity
	}
	if o.AmazonMailSettings.CharSet == "" {
		o.AmazonMailSettings.CharSet = AmazonMailSettingsDefaultCharSet
	}
	if o.EncryptSettings.Salt == "" {
		o.EncryptSettings.Salt = EncryptSettingsDefaultSalt
	}
}
