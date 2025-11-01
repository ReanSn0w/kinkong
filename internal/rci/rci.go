package rci

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

func New(baseURL, cName, cValue string) *Client {
	return &Client{
		baseURL: baseURL,
		cName:   cName,
		cValue:  cValue,
		cl: &http.Client{
			Timeout: time.Second * 3,
		},
	}
}

type Client struct {
	cl      *http.Client
	baseURL string
	cName   string
	cValue  string
}

type InterfaceResponse map[string]InterfaceItem

type InterfaceItem struct {
	Description string `json:"description"`
}

func (c Client) GetWireguardInterfaces() (InterfaceResponse, error) {
	resp, err := c.GetInterfaces()
	if err != nil {
		return nil, err
	}

	for key, _ := range resp {
		if !strings.Contains(key, "Wireguard") {
			delete(resp, key)
		}
	}

	return resp, nil
}

func (c *Client) GetInterfaces() (InterfaceResponse, error) {
	resp, err := c.request(http.MethodGet, "/interface", nil)
	if err != nil {
		return nil, err
	}

	data := InterfaceResponse{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data, err
}

type RouteItem struct {
	Host      string `json:"host"`
	Network   string `json:"network"`
	Mask      string `json:"mask"`
	Interface string `json:"interface"`
	Comment   string `json:"comment"`
}

func (r *RouteItem) Hash() string {
	seed := fmt.Sprintf("%s:%s:%s", r.Network, r.Mask, r.Interface)
	hash := md5.Sum([]byte(seed))
	return fmt.Sprintf("%x", hash)
}

func (r *RouteItem) IsAppRoute() bool {
	return strings.Contains(r.Comment, r.Hash())
}

func (r *RouteItem) GetNetwork() string {
	if r.Host != "" {
		return r.Host
	}

	return r.Network
}

func (r *RouteItem) GetMask() string {
	if r.Host != "" {
		return "255.255.255.255"
	}

	return r.Mask
}

func (c *Client) ListIPRoutes() ([]RouteItem, error) {
	resp, err := c.request(http.MethodGet, "/ip/route", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	data := []RouteItem{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data, err
}

func (c *Client) AddIPRoute(item RouteItem) error {
	_, err := c.request(http.MethodPost, "/ip/route", map[string]interface{}{
		"network":   item.Network,
		"mask":      item.Mask,
		"interface": item.Interface,
		"comment":   item.Comment,
	})
	return err
}

func (c *Client) DeleteIPRoute(item RouteItem) error {
	if item.Mask == "255.255.255.255" {
		_, err := c.request(http.MethodDelete, "/ip/route", map[string]interface{}{
			"host":      strings.Replace(item.Network, "/32", "", -1),
			"interface": item.Interface,
		})

		return err
	}

	_, err := c.request(http.MethodDelete, "/ip/route", map[string]interface{}{
		"network":   item.Network,
		"mask":      item.Mask,
		"interface": item.Interface,
	})
	return err
}

func (c *Client) request(method string, path string, params map[string]interface{}) (*http.Response, error) {
	requestURL := c.baseURL + path
	var body io.Reader
	switch method {
	case http.MethodGet, http.MethodDelete:
		vals := url.Values{}
		for key, val := range params {
			vals.Add(key, fmt.Sprint(val))
		}
		if len(vals) > 0 {
			requestURL += "?" + vals.Encode()
		}
	case http.MethodPost:
		buffer := new(bytes.Buffer)
		err := json.NewEncoder(buffer).Encode(params)
		if err != nil {
			return nil, err
		}
		body = buffer
	}

	request, err := http.NewRequest(method, requestURL, body)
	if err != nil {
		return nil, err
	}

	request.AddCookie(&http.Cookie{
		Name:  c.cName,
		Value: c.cValue,
	})

	resp, err := c.cl.Do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		data, _ := httputil.DumpResponse(resp, true)
		return nil, fmt.Errorf("wrong status code: \n%v", string(data))
	}

	return resp, err
}
