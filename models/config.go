package models

import (
	"time"
)

const (
	CONN_SECURITY_NONE                 = ""
	CONN_SECURITY_PLAIN                = "PLAIN"
	CONN_SECURITY_TLS                  = "TLS"
	CONN_SECURITY_STARTTLS             = "STARTTLS"
	APP_SETTINGS_DEFAULT_PATH          = "http://testtest.twreporter.org:8080"
	APP_SETTINGS_DEFAULT_EXPIRATION    = 168 // 7 days
	EMAIL_SETTINGS_DEFAULT_SMTP_SERVER = "smtp.office365.com"
	EMAIL_SETTINGS_DEFAULT_SMTP_PORT   = "587"
	DB_SETTINGS_DEFAULT_NAME           = "test_membership"
	DB_SETTINGS_DEFAULT_ADDRESS        = "127.0.0.1"
	DB_SETTINGS_DEFAULT_PORT           = "3306"
	ENCRYPT_SETTINGS_DEFAULT_SALT      = "@#$%"
)

type AppSettings struct {
	Path       string
	Token      string
	Expiration time.Duration
}

type EmailSettings struct {
	SMTPUsername       string
	SMTPPassword       string
	SMTPServer         string
	SMTPPort           string
	ConnectionSecurity string
	FeedbackName       string
	FeedbackEmail      string
}

type DBSettings struct {
	Name     string
	User     string
	Password string
	Address  string
	Port     string
}

type FacebookSettings struct {
	Id       string
	Secret   string
	Url      string
	Statestr string
}

type OauthSettings struct {
	FacebookSettings FacebookSettings
}

type EncryptSettings struct {
	Salt string
}

type Config struct {
	AppSettings     AppSettings
	EmailSettings   EmailSettings
	DBSettings      DBSettings
	OauthSettings   OauthSettings
	EncryptSettings EncryptSettings
}

func (o *Config) SetDefaults() {
	if o.AppSettings.Expiration == 0 {
		o.AppSettings.Expiration = APP_SETTINGS_DEFAULT_EXPIRATION
	}

	if o.DBSettings.Name == "" {
		o.DBSettings.Name = DB_SETTINGS_DEFAULT_NAME
	}
	if o.DBSettings.Address == "" {
		o.DBSettings.Address = DB_SETTINGS_DEFAULT_ADDRESS
	}
	if o.DBSettings.Port == "" {
		o.DBSettings.Port = DB_SETTINGS_DEFAULT_PORT
	}
	if o.EmailSettings.SMTPServer == "" {
		o.EmailSettings.SMTPServer = EMAIL_SETTINGS_DEFAULT_SMTP_SERVER
	}
	if o.EmailSettings.SMTPPort == "" {
		o.EmailSettings.SMTPPort = EMAIL_SETTINGS_DEFAULT_SMTP_PORT
	}
	if o.EmailSettings.ConnectionSecurity == "" {
		o.EmailSettings.ConnectionSecurity = CONN_SECURITY_STARTTLS
	}
	if o.EncryptSettings.Salt == "" {
		o.EncryptSettings.Salt = ENCRYPT_SETTINGS_DEFAULT_SALT
	}
}
