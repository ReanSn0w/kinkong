package rci

import "errors"

type Config struct {
	Host        string `long:"host" env:"HOST" default:"192.168.1.1" description:"RCI host"`
	CookieName  string `long:"cookie-name" env:"COOKIE_NAME" description:"RCI cookie name"`
	CookieValue string `long:"cookie-value" env:"COOKIE_VALUE" description:"RCI cookie value"`
}

func (c Config) MustInitRCI() *Client {
	client, err := c.InitRCI()
	if err != nil {
		panic("init rci failed: " + err.Error())
	}
	return client
}

func (c Config) InitRCI() (*Client, error) {
	if c.CookieName == "" || c.CookieValue == "" {
		return nil, errors.New("rci cookie name and value are required")
	}

	client := New("http://"+c.Host+"/rci", c.CookieName, c.CookieValue)
	return client, nil
}
