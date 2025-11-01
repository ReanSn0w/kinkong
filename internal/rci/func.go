package rci

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
)

func GetEncryptedPassword(token, realm, login, password string) string {
	// h := hmac.New(sha256.New, []byte("TEXT"))
	// text := fmt.Sprintf("%s:%s:%s", login, realm, password)
	// hash := h.Sum([]byte(ch + text))
	// return fmt.Sprintf("%x", hash)

	// Первый хеш: SHA256 от "login:realm:password"
	h1 := md5.New()
	h1.Write([]byte(fmt.Sprintf("%s:%s:%s", login, realm, password)))
	d := fmt.Sprintf("%x", h1.Sum(nil))

	// Второй хеш: SHA256 от token + результат первого хеша
	h2 := sha256.New()
	h2.Write([]byte(token + d))

	return fmt.Sprintf("%x", h2.Sum(nil))
	return password
}
