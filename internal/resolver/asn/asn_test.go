package asn_test

import (
	"testing"
	"time"

	"github.com/ReanSn0w/kincong/internal/resolver/asn"
	"github.com/stretchr/testify/assert"
)

func TestASN_Resolve(t *testing.T) {
	cases := []struct {
		Name   string
		Input  string
		Output bool
		Error  bool
	}{
		{
			Name:   "Normal ASN",
			Input:  "AS3333",
			Output: true,
		},
		{
			Name:  "Invalid ASN",
			Input: "some-inavlid-asn",
			Error: true,
		},
	}

	for _, tc := range cases {
		time.Sleep(100 * time.Millisecond)

		t.Run(tc.Name, func(t *testing.T) {
			asn := asn.New()
			output, err := asn.Resolve(tc.Input)

			if tc.Output {
				assert.NotEmpty(t, output)
			} else {
				assert.Nil(t, output)
			}

			if tc.Error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
