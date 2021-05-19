package distkvs

import (
	"errors"
	"log"

	"example.org/cpsc416/a5/kvslib"
	"github.com/DistributedClocks/tracing"
	// "log"
)

const ChCapacity = 10

type ClientConfig struct {
	ClientID         string
	FrontEndAddr     string
	TracerServerAddr string
	TracerSecret     []byte
}

type Client struct {
	NotifyChannel kvslib.NotifyChannel
	id            string
	frontEndAddr  string
	kvs           *kvslib.KVS
	tracer        *tracing.Tracer
	initialized   bool
	tracerConfig  tracing.TracerConfig
}

func NewClient(config ClientConfig, kvs *kvslib.KVS) *Client {
	t_config := tracing.TracerConfig{config.TracerServerAddr, config.ClientID, config.TracerSecret}
	tracer := tracing.NewTracer(t_config)
	// log.Printf("k: %+v\n", k)

	c := Client{
		NotifyChannel: nil,
		id:            config.ClientID,
		frontEndAddr:  config.FrontEndAddr,
		tracer:        tracer,
		initialized:   false,
		tracerConfig:  t_config,
		kvs:           kvs}
	// log.Printf("%s",c.id)
	return &c
}

func (c *Client) Initialize() error {
	// Call KVS initialize here
	notifyCh, err := c.kvs.Initialize(c.tracer, c.id, c.frontEndAddr, ChCapacity)
	if err != nil {
		return errors.New("kvs initialize error")
	}
	c.NotifyChannel = notifyCh

	if err == nil {
		c.initialized = true
	}

	if c.initialized == true {
		return nil
	}
	log.Print(err)

	return errors.New("Client Cannot be initialized")
}

func (c *Client) Get(clientId string, key string) (uint32, error) {
	return c.kvs.Get(c.tracer, clientId, key)
}

func (c *Client) Put(clientId string, key string, value string, delay int) (uint32, error) {
	return c.kvs.Put(c.tracer, clientId, key, value, delay)
}

func (c *Client) Close() error {
	return c.kvs.Close()
}
