package models

import (
	"time"
)

const (
	// Office360 global string constants
	Office360 = "Office360"

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
)

// AppSettings could be defined in configs/config.json
type AppSettings struct {
	Path       string
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

// DBSettings could be defined in configs/config.json
type DBSettings struct {
	Name     string
	User     string
	Password string
	Address  string
	Port     string
}

// FacebookSettings could be defined in configs/config.json
type FacebookSettings struct {
	ID       string
	Secret   string
	URL      string
	Statestr string
}

// OauthSettings could be defined in configs/config.json
type OauthSettings struct {
	FacebookSettings FacebookSettings
}

// EncryptSettings could be defined in configs/config.json
type EncryptSettings struct {
	Salt string
}

// Config contains all the other configs
type Config struct {
	AppSettings     AppSettings
	EmailSettings   EmailSettings
	DBSettings      DBSettings
	OauthSettings   OauthSettings
	EncryptSettings EncryptSettings
}

// SetDefaults could set default value in the Config struct
func (o *Config) SetDefaults() {
	if o.AppSettings.Expiration == 0 {
		o.AppSettings.Expiration = AppSettingsDefaultExpiration
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
	if o.EncryptSettings.Salt == "" {
		o.EncryptSettings.Salt = EncryptSettingsDefaultSalt
	}
}
