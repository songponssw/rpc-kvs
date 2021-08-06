package distkvs

import (
	"errors"
	"log"

	"example.org/cpsc416/a5/kvslib"
	// "log"
)

const ChCapacity = 10

type ClientConfig struct {
	ClientID     string
	FrontEndAddr string
}

type Client struct {
	NotifyChannel kvslib.NotifyChannel
	id            string
	frontEndAddr  string
	kvs           *kvslib.KVS
	initialized   bool
}

func NewClient(config ClientConfig, kvs *kvslib.KVS) *Client {
	// log.Printf("k: %+v\n", k)

	c := Client{
		NotifyChannel: nil,
		id:            config.ClientID,
		frontEndAddr:  config.FrontEndAddr,
		initialized:   false,
		kvs:           kvs}
	// log.Printf("%s",c.id)
	return &c
}

func (c *Client) Initialize() error {
	// Call KVS initialize here
	notifyCh, err := c.kvs.Initialize(c.id, c.frontEndAddr, ChCapacity)
	if err != nil {
		return err
		// return errors.New("kvs initialize error")
	}
	c.NotifyChannel = notifyCh

	if err != nil {
		log.Fatal(err)
		return errors.New("Client Cannot be initialized")
	}

	c.initialized = true

	return nil
}

func (c *Client) Get(reqId uint32, key string) (uint32, error) {
	return c.kvs.Get(reqId, key)
}

func (c *Client) Put(reqId uint32, key string, value string, delay int) (uint32, error) {
	return c.kvs.Put(reqId, key, value, delay)
}

func (c *Client) Close() error {
	return c.kvs.Close()
}
