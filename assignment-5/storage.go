package distkvs

import (
	"errors"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/DistributedClocks/tracing"
)

type StorageConfig struct {
	StorageID        string
	StorageAdd       StorageAddr
	ListenAddr       string
	FrontEndAddr     string
	DiskPath         string
	TracerServerAddr string
	TracerSecret     []byte
}

type StorageLoadSuccess struct {
	State map[string]string
}

type StoragePut struct {
	Key   string
	Value string
}

type StorageSaveData struct {
	Key   string
	Value string
}

type StorageGet struct {
	Key string
}

type StorageGetResult struct {
	Key   string
	Value *string
}

type Storage struct {
	// state may go here
}

func (*Storage) Start(frontEndAddr string, storageAddr string, diskPath string, strace *tracing.Tracer) error {

	result := new(Storage)
	rpc.Register(result)

	rpc.HandleHTTP()

	l, e := net.Listen("tcp", frontEndAddr)
	if e != nil {
		log.Fatal("listen error:", e)
	}

	err := http.Serve(l, nil)
	if err != nil {
		log.Fatal("listen error:", err)
	}

	return errors.New("not implemented")
}

func (*Storage) StorageGet(args StorageGet, reply *StorageGetResult) error {
	ret := StorageGetResult{}
	s := args.Key + "Value"
	ret.Value = &s
	*reply = ret
	log.Printf("Get value %s from %s", *reply.Value, args.Key)
	return nil
}

func (*Storage) StoragePut(args StoragePut, reply *string) error {
	log.Printf("Put value %s to %s", args.Value, args.Key)
	s := "Success"
	// sPtr := new(string)
	// sPtr = &s
	*reply = s

	log.Print(*reply)
	return nil
}
