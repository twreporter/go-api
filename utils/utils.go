package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/scrypt"
)

// GenerateActivateMailBody generate the html a tag which can link to /active enpoint to activate the account
func GenerateActivateMailBody(mailAddress, activeToken string) string {
	href := fmt.Sprintf("%s://%s:%s/activate?email=%s&token=%s", Cfg.ConsumerSettings.Protocal, Cfg.ConsumerSettings.Host, Cfg.ConsumerSettings.Port, mailAddress, activeToken)

	// TBD make the activate mail more beautiful and informative
	return fmt.Sprintf("<a href=\"%s\" target=\"_blank\">Activate Your Account</a>", href)
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// GenerateEncryptedPassword returns encryptedly
// securely generated string.
func GenerateEncryptedPassword(password []byte) (string, error) {
	salt := []byte(Cfg.EncryptSettings.Salt)
	key, err := scrypt.Key(password, salt, 16384, 8, 1, 32)
	return fmt.Sprintf("%x", key), err
}
