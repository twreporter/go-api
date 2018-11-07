package configs

import (
	"bytes"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

var defaultConf = []byte(`
environment: development
cors:
    allow_origins:
        - 'http://localhost:3000'
        - 'http://localhost:3001'
app:
    protocol: http
    host: localhost
    port: '8080'
    domain: localhost
    jwt_secret: secret_token
    jwt_expiration: 604800
    jwt_issuer: 'http://testtest.twreporter.org:8080' # used for issuer claim
    jwt_audience: 'http://testtest.twreporter.org:3000' # used for audience claim
email:
    smtp:
        username: no-reply@t-reporters.org
        password: smtp_password
        server: smtp.office365.com
        server_owner: office365
        port: '587'
        connection_security: STARTTLS
        feedback_name: '報導者 The Reporter'
        feedback_email: contact@twreporter.org
    amazon:
        sender: no-reply@twreporter.org
        aws_region: us-west-2
        char_set: utf-8
db:
    mysql:
        name: test_membership
        user: test_membership
        password: test_membership
        address: 127.0.0.1
        port: '3306'
    mongo:
        url: localhost
        dbname: gorm
        timeout: 5
oauth:
    facebook:
        id: "" # provide your own facebook oauth ID
        secret: "" # provide your own facebook oauth secret
    google:
        id: "" # provide your own ID
        secret: "" # provide your own secret
donation:
    card_secret_key: test_card_secret_key
    tappay_url: 'https://sandbox.tappaysdk.com/tpc/payment/pay-by-prime'
    tappay_partner_key: 'partner_6ID1DoDlaPrfHw6HBZsULfTYtDmWs0q0ZZGKMBpp4YICWBxgK97eK3RM'
algolia:
    application_id: "" # provide your own application ID
    api_key: "" # provide your own api key
encrypt:
    salt: '@#$%'
`)

type ConfYaml struct {
	Environment string         `yaml:"environment"`
	Cors        CorsConfig     `yaml:"cors"`
	App         AppConfig      `yaml:"app"`
	Email       EmailConfig    `yaml:"email"`
	DB          DBConfig       `yaml:"db"`
	Oauth       OauthConfig    `yaml:"oauth"`
	Donation    DonationConfig `yaml:"donation"`
	Algolia     AlgoliaConfig  `ymal:"algolia"`
	Encrypt     EncryptConfig  `yaml:"encrypt"`
}

type CorsConfig struct {
	AllowOrigins []string `yaml:"allow_origins"`
}

type AppConfig struct {
	Protocol      string `yaml:"protocol"`
	Host          string `yaml:"host"`
	Port          string `yaml:"port"`
	Domain        string `yaml:"domain"`
	JwtSecret     string `yaml:"jwt_secret"`
	JwtExpiration int    `yaml:"jwt_expiration"`
	JwtIssuer     string `yaml:"jwt_issuer"`
	JwtAudience   string `yaml:"jwt_audience"`
}

type EmailConfig struct {
	SMTP   SMTPConfig   `yaml:"smtp"`
	Amazon AmazonConfig `yaml:"amazon"`
}

type SMTPConfig struct {
	Username           string `yaml:"username"`
	Password           string `yaml:"password"`
	Server             string `yaml:"server"`
	ServerOwner        string `yaml:"server_owner"`
	Port               string `yaml:"port"`
	ConnectionSecurity string `yaml:"connection_security"`
	FeedbackName       string `yaml:"feedback_name"`
	FeedbackEmail      string `yaml:"feedback_email"`
}

type AmazonConfig struct {
	Sender    string `yaml:"sender"`
	AwsRegion string `yaml:"aws_region"`
	Charset   string `yaml:"char_set"`
}

type DBConfig struct {
	MySQL MySQLConfig `yaml:"mysql"`
	Mongo MongoConfig `yaml:"mongo"`
}

type MySQLConfig struct {
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Address  string `yaml:"address"`
	Port     string `yaml:"port"`
}

type MongoConfig struct {
	URL     string `yaml:"url"`
	DBname  string `yaml:"dbname"`
	Timeout int    `yaml:"timeout"`
}

type OauthConfig struct {
	Facebook FacebookConfig `yaml:"facebook"`
	Google   GoogleConfig   `yaml:"google"`
}

type FacebookConfig struct {
	ID     string `yaml:"id"`
	Secret string `yaml:"secret"`
}

type GoogleConfig struct {
	ID     string `yaml:"id"`
	Secret string `yaml:"secret"`
}

type DonationConfig struct {
	CardSecretKey    string `yaml:"card_secret_key"`
	TapPayURL        string `yaml:"tappay_url"`
	TapPayPartnerKey string `yaml:"tappay_partner_key"`
}

type AlgoliaConfig struct {
	ApplicationID string `yaml:"application_id"`
	APIKey        string `yaml:"api_key"`
}

type EncryptConfig struct {
	Salt string `yaml:"salt"`
}

func init() {
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()        // read in environment variables that match
	viper.SetEnvPrefix("goapi") // will be uppercased automatically
}

func buildConf() ConfYaml {
	var conf ConfYaml

	// Environemt
	conf.Environment = viper.GetString("environment")

	// App
	conf.App.Host = viper.GetString("app.host")
	conf.App.Protocol = viper.GetString("app.protocol")
	conf.App.Port = viper.GetString("app.port")
	conf.App.Domain = viper.GetString("app.domain")
	conf.App.JwtSecret = viper.GetString("app.jwt_secret")
	conf.App.JwtExpiration = viper.GetInt("app.jwt_expiration")
	conf.App.JwtAudience = viper.GetString("app.jwt_audience")
	conf.App.JwtIssuer = viper.GetString("app.jwt_issuer")

	// Cors
	conf.Cors.AllowOrigins = viper.GetStringSlice("cors.allow_origins")

	// DB - MySQL
	conf.DB.MySQL.Name = viper.GetString("db.mysql.name")
	conf.DB.MySQL.Password = viper.GetString("db.mysql.password")
	conf.DB.MySQL.Address = viper.GetString("db.mysql.address")
	conf.DB.MySQL.Port = viper.GetString("db.mysql.port")
	conf.DB.MySQL.User = viper.GetString("db.mysql.user")

	// DB - Mongo
	conf.DB.Mongo.DBname = viper.GetString("db.mongo.dbname")
	conf.DB.Mongo.URL = viper.GetString("db.mongo.url")
	conf.DB.Mongo.Timeout = viper.GetInt("db.mongo.timeout")

	// Email - Amazon
	conf.Email.Amazon.Sender = viper.GetString("email.amazon.sender")
	conf.Email.Amazon.AwsRegion = viper.GetString("email.amazon.aws_region")
	conf.Email.Amazon.Charset = viper.GetString("email.amazon.char_set")

	// Email - SMTP
	conf.Email.SMTP.Server = viper.GetString("email.smtp.server")
	conf.Email.SMTP.ServerOwner = viper.GetString("email.smtp.server_owner")
	conf.Email.SMTP.ConnectionSecurity = viper.GetString("email.smtp.connection_security")
	conf.Email.SMTP.Port = viper.GetString("email.smtp.port")
	conf.Email.SMTP.Username = viper.GetString("email.smtp.username")
	conf.Email.SMTP.Password = viper.GetString("email.smtp.password")
	conf.Email.SMTP.FeedbackEmail = viper.GetString("email.smtp.feedback_email")
	conf.Email.SMTP.FeedbackName = viper.GetString("email.smtp.feedback_name")

	// Oauth - Facebook
	conf.Oauth.Facebook.ID = viper.GetString("oauth.facebook.id")
	conf.Oauth.Facebook.Secret = viper.GetString("oauth.facebook.secret")

	// Oauth - Google
	conf.Oauth.Google.ID = viper.GetString("oauth.google.id")
	conf.Oauth.Google.Secret = viper.GetString("oauth.google.secret")

	// TapPay
	conf.Donation.CardSecretKey = viper.GetString("donation.card_secret_key")
	conf.Donation.TapPayURL = viper.GetString("donation.tappay_url")
	conf.Donation.TapPayPartnerKey = viper.GetString("donation.tappay_partner_key")

	// Algolia
	conf.Algolia.ApplicationID = viper.GetString("algolia.application_id")
	conf.Algolia.APIKey = viper.GetString("algolia.api_key")

	// Encrypt
	conf.Encrypt.Salt = viper.GetString("encrypt.salt")

	return conf
}

// LoadDefaultConf loads default config
func LoadDefaultConf() (ConfYaml, error) {
	var conf ConfYaml

	// load default config
	if err := viper.ReadConfig(bytes.NewBuffer(defaultConf)); err != nil {
		return conf, err
	}

	conf = buildConf()

	return conf, nil
}

// LoadConf load config from file and read in environment variables that match
func LoadConf(confPath string) (ConfYaml, error) {
	var conf ConfYaml

	if confPath != "" {
		content, err := ioutil.ReadFile(confPath)

		if err != nil {
			return conf, err
		}

		if err := viper.ReadConfig(bytes.NewBuffer(content)); err != nil {
			return conf, err
		}
	} else {
		viper.AddConfigPath("$GOPATH/src/twreporter.org/go-api/configs/")
		viper.AddConfigPath("./configs/")
		viper.SetConfigName("config")

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			log.Infof("Using config file: %s", viper.ConfigFileUsed())
		} else {
			// load default config
			return LoadDefaultConf()
		}
	}

	conf = buildConf()

	return conf, nil
}
