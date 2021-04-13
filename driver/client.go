package driver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Addr   string
	Port   int
	urlSet string
	urlGet string
}

type KVPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewClient(addr string, port int) *Client {
	client := &Client{
		Addr: addr,
		Port: port,
	}

	client.urlSet = fmt.Sprintf("http://%s:%d/store/set", client.Addr, client.Port)
	client.urlGet = fmt.Sprintf("http://%s:%d/store/get/", client.Addr, client.Port)
	return client
}

func (c *Client) Set(key, value string) error {
	kvp := &KVPair{
		Key:   key,
		Value: value,
	}
	bs, err := json.Marshal(kvp)
	if err != nil {
		return err
	}
	_, err = http.Post(c.urlSet, "application/json", bytes.NewReader(bs))
	return err
}

func (c *Client) Get(key string) (string, error) {
	resp, err := http.Get(c.urlGet + key)
	if err != nil {
		return "", err
	}

	value, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(value), nil
}
