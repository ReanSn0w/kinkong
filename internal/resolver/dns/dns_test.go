package dns_test

import (
	"sort"
	"testing"

	"github.com/ReanSn0w/kincong/internal/resolver/dns"
	"github.com/ReanSn0w/kincong/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestDNS_Resolve(t *testing.T) {
	cases := []struct {
		Name   string
		Input  string
		Output []utils.ResolvedSubnet
	}{
		{
			Name:  "Normal Domain",
			Input: "coder.papkovda.ru",
			Output: []utils.ResolvedSubnet{
				"212.34.135.116/32",
			},
		},
		{
			Name:   "Invalid Domain",
			Input:  "invalid-domain",
			Output: nil,
		},
	}

	resolver := dns.New("8.8.8.8")

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			output, _ := resolver.Resolve(tc.Input)

			sort.Slice(tc.Output, func(i, j int) bool {
				return tc.Output[i] < tc.Output[j]
			})

			sort.Slice(output, func(i, j int) bool {
				return output[i] < output[j]
			})

			assert.Equal(t, tc.Output, output)
		})
	}
}
