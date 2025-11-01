package rci

import "errors"

type Config struct {
	Host     string `long:"host" env:"HOST" default:"192.168.1.1" description:"RCI host"`
	Login    string `long:"login" env:"LOGIN" description:"router user login"`
	Password string `long:"password" env:"PASSWORD" description:"router user password"`
}

func (c Config) MustInitRCI() *Client {
	client, err := c.InitRCI()
	if err != nil {
		panic("init rci failed: " + err.Error())
	}
	return client
}

func (c Config) InitRCI() (*Client, error) {
	if c.Login == "" || c.Password == "" {
		return nil, errors.New("rci login and password are required")
	}

	return New("http://"+c.Host, c.Login, c.Password)
}
