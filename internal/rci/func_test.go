package rci_test

import (
	"testing"

	"github.com/ReanSn0w/kincong/internal/rci"
	"github.com/stretchr/testify/assert"
)

func Test_GetEncryptedPassword(t *testing.T) {
	cases := []struct {
		Name      string
		Challenge string
		Realm     string
		Login     string
		Password  string
		Hash      string
	}{
		{
			Name:      "Test case 1",
			Challenge: "XIQGOHVLSYRSQKSIFDDIRQDAPOEHCGFL",
			Realm:     "Keenetic Giga",
			Login:     "config",
			Password:  "123412341234",
			Hash:      "e225695d4f080672fe35c0e474d072e9ee784cb2900a75a5caeda82071e93a2b",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			hash := rci.GetEncryptedPassword(tc.Challenge, tc.Realm, tc.Login, tc.Password)
			assert.Equal(t, tc.Hash, hash)
		})
	}

}
