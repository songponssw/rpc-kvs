package distkvs

import (
	"errors"

	"example.org/cpsc416/a5/kvslib"
	"github.com/DistributedClocks/tracing"

//	"log"
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
	c := Client{
		id: config.ClientID,
		frontEndAddr: config.FrontEndAddr,
		kvs: kvs,
		tracer: tracing.NewTracer(t_config),
		initialized: false,
		tracerConfig: t_config}
	return &c
}

func (c *Client) Initialize() error {
	c.initialized = true
	return nil
	return errors.New("Client Init not implemented")
}

func (c *Client) Get(clientId string, key string) (uint32, error) {
	return c.kvs.Get(c.tracer, clientId, key)
}

func (c *Client) Put(clientId string, key string, value string) (uint32, error) {
	return c.kvs.Put(c.tracer, clientId, key, value)
}

func (c *Client) Close() error {
	return c.kvs.Close()
}
