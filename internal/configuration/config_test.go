package configuration_test

import (
	"os"
	"strings"
	"testing"

	"github.com/ReanSn0w/kincong/internal/configuration"
	"github.com/ReanSn0w/kincong/internal/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

const (
	filename = "config.yaml"
)

func Test_Load(t *testing.T) {
	cases := []struct {
		Name     string
		Config   []byte
		Expected configuration.Configuration
	}{
		{
			Name: "Load Valid Configuration",
			Config: []byte(strings.ReplaceAll(`
				- title: Test Configuration
				  values:
				    - 192.168.0.1
				    - http://mysite.org
				    - AS13335
			`, "\t", "")),
			Expected: configuration.Configuration{
				{
					Title: "Test Configuration",
					Values: []utils.Value{
						"192.168.0.1",
						"http://mysite.org",
						"AS13335",
					},
				},
			},
		},
		{
			Name: "Load Valid Configuration With Another Fields",
			Config: []byte(strings.ReplaceAll(`
				- title: Test Configuration
				  description: This is a test configuration
				  values:
				    - 192.168.0.1
				    - http://mysite.org
				    - AS13335
			`, "\t", "")),
			Expected: configuration.Configuration{
				{
					Title: "Test Configuration",
					Values: []utils.Value{
						"192.168.0.1",
						"http://mysite.org",
						"AS13335",
					},
				},
			},
		},
	}

	tmpDir := os.TempDir()
	filepath := tmpDir + "/" + filename
	defer os.Remove(filepath)

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			func() {
				f, err := os.Create(filepath)
				if err != nil {
					t.Fatal(err)
				}

				defer f.Close()

				_, err = f.Write(tc.Config)
				if err != nil {
					t.Fatal(err)
				}
			}()

			func() {
				config, err := configuration.Load(filepath)
				assert.NoError(t, err)
				if assert.NotNil(t, config) {
					assert.Equal(t, tc.Expected, config)
				}
			}()
		})
	}
}

func Test_Configuration(t *testing.T) {
	exampleConfiguration := configuration.Configuration{
		{
			Title: "Test Configuration",
			Values: []utils.Value{
				"192.168.0.1",
				"http://mysite.org",
				"AS13335",
			},
		},
	}

	f, err := os.Create("./" + filename)
	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()

	err = yaml.NewEncoder(f).Encode(exampleConfiguration)
	if err != nil {
		t.Fatal(err)
	}
}
