package dns

type Config struct {
	Enabled bool     `long:"enabled" env:"ENABLED" description:"enable DNS resolver"`
	Host    []string `long:"host" env:"HOST" default:"8.8.8.8" description:"DNS host"`
}

func (c Config) InitDNSResolver() *Resolver {
	if !c.Enabled {
		return nil
	}
	return New(c.Host...)
}
