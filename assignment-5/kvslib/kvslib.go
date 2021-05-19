// Package kvslib provides an API which is a wrapper around RPC calls to the
// frontend.
package kvslib

import (
	"errors"
	"log"

	"github.com/DistributedClocks/tracing"

	"net/rpc"
)

type KvslibBegin struct {
	ClientId string
}

type KvslibPut struct {
	ClientId string
	OpId     uint32
	Key      string
	Value    string
}

type KvslibGet struct {
	ClientId string
	OpId     uint32
	Key      string
}

type KvslibPutResult struct {
	OpId uint32
	Err  bool
}

type KvslibGetResult struct {
	OpId  uint32
	Key   string
	Value *string
	Err   bool
}

type KvslibComplete struct {
	ClientId string
}

// NotifyChannel is used for notifying the client about a mining result.
type NotifyChannel chan ResultStruct

type ResultStruct struct {
	OpId        uint32
	StorageFail bool
	Result      *string
}

type KVS struct {
	notifyCh  NotifyChannel
	rpcClient *rpc.Client
	OpId      uint32
	ClientId  KvslibBegin
	// Add more KVS instance state here.
}

func NewKVS() *KVS {
	return &KVS{
		notifyCh: nil,
	}
}

// Initialize Initializes the instance of KVS to use for connecting to the frontend,
// and the frontends IP:port. The returned notify-channel channel must
// have capacity ChCapacity and must be used by kvslib to deliver all solution
// notifications. If there is an issue with connecting, this should return
// an appropriate err value, otherwise err should be set to nil.
func (d *KVS) Initialize(localTracer *tracing.Tracer, clientId string, frontEndAddr string, chCapacity uint) (NotifyChannel, error) {
	// dial
	rpcClient, err := rpc.DialHTTP("tcp", frontEndAddr)
	if err != nil {
		return nil, errors.New("Cannot established connection with RPC server.")
	}
	d.rpcClient = rpcClient
	d.OpId = 0

	d.ClientId = KvslibBegin{clientId}

	notifyLocal := make(chan ResultStruct, chCapacity)
	d.notifyCh = notifyLocal

	return d.notifyCh, nil
}

// Get is a non-blocking request from the client to the system. This call is used by
// the client when it wants to get value for a key.
func (d *KVS) Get(tracer *tracing.Tracer, clientId string, key string) (uint32, error) {
	d.OpId++
	args := KvslibGet{d.ClientId.ClientId, d.OpId, key}
	reply := new(ResultStruct) // This shoulbe GetResult Struct???

	funcCall := d.rpcClient.Go("FrontEnd.HandleGet", args, &reply, nil)
	replyCall := <-funcCall.Done

	// Log result using Trancer???
	log.Print(*reply.Result)
	reply.OpId = d.OpId
	d.notifyCh <- *reply
	log.Print(d.notifyCh)

	log.Print("added to channel")

	if replyCall.Error != nil {
		return d.OpId, errors.New("key not found")
	}

	// Should return OpId or error
	return d.OpId, nil
}

// Put is a non-blocking request from the client to the system. This call is used by
// the client when it wants to update the value of an existing key or add add a new
// key and value pair.
func (d *KVS) Put(tracer *tracing.Tracer, clientId string, key string, value string) (uint32, error) {
	d.OpId += 1
	args := KvslibPut{d.ClientId.ClientId, d.OpId, key, value}

	reply := new(ResultStruct)
	funcCall := d.rpcClient.Go("FrontEnd.HandlePut", args, &reply, nil)
	replyCall := <-funcCall.Done

	log.Print(*reply.Result)

	reply.OpId = d.OpId
	d.notifyCh <- *reply
	log.Print(d.notifyCh)

	//Hanle key not Fond : storage will create the new
	if replyCall != nil {
		return d.OpId, errors.New("key not found")
	}

	d.notifyCh <- *reply
	//Handle update key
	return d.OpId, nil
}

// Close Stops the KVS instance from communicating with the frontend and
// from delivering any solutions via the notify-channel. If there is an issue
// with stopping, this should return an appropriate err value, otherwise err
// should be set to nil.
func (d *KVS) Close() error {
	err := d.rpcClient.Close()
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}
