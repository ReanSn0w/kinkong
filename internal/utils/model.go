package utils

import (
	"crypto/md5"
	"fmt"
	"net"
	"strings"
)

// Value - Значение введенное пользоватем
type Value string

// IsSite -	Проверка на то, что значение является сайтом
func (v Value) IsSite() bool {
	return strings.HasPrefix(string(v), "http")
}

// IsASN - Проверка на то, что значение является ASN
func (v Value) IsASN() bool {
	return strings.HasPrefix(string(v), "AS")
}

// IsIP - Проверка на то, что значение является IP
func (v Value) IsIP() bool {
	return net.ParseIP(string(v)) != nil
}

// IsNetwork - Проверка на то, что значение является сетью
func (v Value) IsNetwork() bool {
	if strings.Contains(string(v), "/") && !v.IsSite() {
		return true
	}

	return false
}

// ResolvedSubnet - Значение добавляемое в конфигурацию роутера
type ResolvedSubnet string

// Hash - Возвращает хеш значения
//
// Планируется использование его в качестве части описания маршрута на роутере
// для определения маршрутов, добавленных с помощью конфигурации, на случае если
// добавленные ранее маршруты придется удалить
func (r ResolvedSubnet) Hash() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(r)))
}
