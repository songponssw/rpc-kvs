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
	k, _ := kvs.Initialize(tracer, config.ClientID, config.FrontEndAddr, ChCapacity)
	log.Printf("k: %+v\n", k)

	c := Client{
		//NotifyChannel ????
		id:           config.ClientID,
		frontEndAddr: config.FrontEndAddr,
		tracer:       tracer,
		initialized:  true,
		tracerConfig: t_config,
		// InitialKVS()
		kvs: kvs}
	// log.Printf("%s",c.id)
	return &c
}

func (c *Client) Initialize() error {
	if c.initialized == true {
		return nil
	}
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
