package asn

type Config struct {
	Enabled bool   `long:"enabled" env:"ENABLED" description:"enable ASN resolver"`
	Key     string `long:"key" env:"KEY" description:"BGPview API key"`
}
