package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"golang.org/x/crypto/scrypt"

	"twreporter.org/go-api/globals"
	//log "github.com/Sirupsen/logrus"
)

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

// GetProjectRoot returns absolute path of current project root.
func GetProjectRoot() string {
	type emptyStruct struct{}
	const rootPkg = "main"

	// use the reflect package to retrieve current package path
	// [go module name]/[package name]
	// i.e. twreporter.org/go-api/utils
	pkg := reflect.TypeOf(emptyStruct{}).PkgPath()
	pkgWithoutModPrefix := string([]byte(pkg)[strings.LastIndex(pkg, "/")+1:])

	// Get the file name of current function
	_, file, _, _ := runtime.Caller(0)

	root := ""
	if pkgWithoutModPrefix == rootPkg {
		// If package lies in main package,
		// then the file must be within project root
		root = string([]byte(file)[:strings.LastIndex(file, "/")])
	} else {
		// If package lies in any sub package,
		// then the project root is the prefix of the current file name.
		root = string([]byte(file)[:strings.Index(file, pkgWithoutModPrefix)-1])
	}

	return root
}
