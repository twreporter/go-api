package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"

	"golang.org/x/crypto/scrypt"

	"twreporter.org/go-api/globals"
	//log "github.com/Sirupsen/logrus"
)

// GenerateActivateMailBody generate the html a tag which can link to /active enpoint to activate the account
func GenerateActivateMailBody(mailAddress, activeToken, destination, host string) string {
	if "" == host {
		host = globals.Conf.App.Host
	}

	href := fmt.Sprintf("%s://%s:%s/activate?email=%s&token=%s&destination=%s", globals.Conf.App.Protocol, host, globals.Conf.App.Port, mailAddress, activeToken, destination)

	var defaultBody = fmt.Sprintf("<a href=\"%s\" target=\"_blank\">啟動報導者帳號</a>", href)

	// activateTpl is from mail_templates.go
	t, err := template.New("webpage").Parse(activateTpl)

	if err != nil {
		return defaultBody
	}

	data := struct {
		Activate string
		Desc     string
		Header   string
		Href     string
	}{
		Activate: "登入《報導者》帳號",
		Desc: `
		親愛的讀者 您好：

		請點擊上方按鈕登入《報導者》網站。

		此連結僅能被使用一次，並且在十五分鐘之後失效。

		若使用上方按鈕無法登入，請複製以下網址，並貼到瀏覽器網址列上。

		` + href,
		Header: "歡迎登入《報導者》，開啟新聞小革命",
		Href:   href,
	}

	var out bytes.Buffer

	err = t.Execute(&out, data)

	if err != nil {
		return defaultBody
	}

	return out.String()
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
	salt := []byte(globals.Conf.Encrypt.Salt)
	key, err := scrypt.Key(password, salt, 16384, 8, 1, 32)
	return fmt.Sprintf("%x", key), err
}
