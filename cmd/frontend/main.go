package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"

	distkvs "example.org/cpsc416/a5"
	"example.org/cpsc416/a5/pb"
	"google.golang.org/grpc"
)

type StorageAddr string

type FrontEndConfig struct {
	ClientAPIListenAddr  string
	StorageAPIListenAddr string
	Storage              StorageAddr
}

type FrontEndStorageFailed struct{}

type FrontEndStorageStarted struct{}

type FrontEndGet struct {
	Key string
}

type FrontEndGetResult struct {
	Key   string
	Value *string
	Err   bool
}

type StorageGet struct {
	Key string
}

type StoragePut struct {
	Key   string
	Value string
	delay int
}

type StorageGetResult struct {
	Key   string
	Value *string
}

type FrontEnd struct {
	// state may go here
	rpcClient *rpc.Client
}

func main() {
	var config FrontEndConfig
	err := distkvs.ReadJSONConfig("config/frontend_config.json", &config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(config)

	frontend := FrontEnd{}
	err = frontend.Start(config.ClientAPIListenAddr, config.StorageAPIListenAddr, 0)

	if err != nil {
		log.Fatal(err)
	}
}

func (f *FrontEnd) Start(clientAPIListenAddr string, storageAPIListenAddr string, storageTimeout uint8) error {
	// result := new(FrontEndGetResult)
	result := new(FrontEnd)
	// rpc.Register(result)

	// result2 := new(FrontEndPutResult)
	// rpc.Register(result2)

	// rpc.HandleHTTP()

	rpcClient, rpcErr := rpc.DialHTTP("tcp", storageAPIListenAddr)
	if rpcErr != nil {
		return errors.New("cannot established connection with Storage server")
	}

	result.rpcClient = rpcClient

	// l, e := net.Listen("tcp", clientAPIListenAddr)
	// if e != nil {
	// 	log.Fatal("listen error:", e)
	// }

	// err := http.Serve(l, nil)
	// if err != nil {
	// 	log.Fatal("listen error:", err)
	// }

	// serve as a grpc server
	grpcServer := grpc.NewServer()
	pb.RegisterFrontendServer(grpcServer, result)

	lis, err := net.Listen("tcp", clientAPIListenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %s", clientAPIListenAddr)
	}

	// start serving
	log.Println("serving on port:", clientAPIListenAddr)
	log.Fatal(grpcServer.Serve(lis))

	return errors.New("frontend Fail")
}

func (f *FrontEnd) HandleGet(
	ctx context.Context, 
	req *pb.FrontendGetRequest,
) (*pb.FrontendGetReponse, error) {

	// convert to storage args
	storageArgs := StorageGet{req.Key}
	storageReply := new(StorageGetResult)
	funcCall := f.rpcClient.Go("Storage.StorageGet", storageArgs, &storageReply, nil)
	replyCall := <-funcCall.Done

	// log.Print(replyCall.Error)
	// log.Print(*storageReply.Value)

	if replyCall.Error != nil {
		return nil, errors.New("FE to Strage fail")
	}

	// ret := kvslib.ResultStruct{}
	// ret.Result = storageReply.Value
	// *reply = ret

	log.Printf("OpId: %d Get value %s from %s", req.OpId, *storageReply.Value, req.Key)

	return &pb.FrontendGetReponse{
		OpId: req.OpId,
		StorageFail: false,
		Result: *storageReply.Value,
	}, nil
}

// !!! Old rpc function
// func (f *FrontEnd) HandleGet(args kvslib.KvslibGet, reply *kvslib.ResultStruct) error {
// 	log.Print("Fuck frontend")
// 	storageArgs := StorageGet{args.Key}
// 	storageReply := new(StorageGetResult)
// 	funcCall := f.rpcClient.Go("Storage.StorageGet", storageArgs, &storageReply, nil)
// 	replyCall := <-funcCall.Done

// 	// log.Print(replyCall.Error)
// 	// log.Print(*storageReply.Value)

// 	if replyCall.Error != nil {
// 		return errors.New("FE to Strage fail")
// 	}

// 	// ret := kvslib.ResultStruct{}
// 	// ret.Result = storageReply.Value
// 	// *reply = ret

// 	// log.Printf("OpId: %d Get value %s from %s", args.OpId, *reply.Result, args.Key)

// 	return nil
// }

func (f *FrontEnd) HandlePut(
	ctx context.Context, 
	req *pb.FrontendPutRequest, 
	) (*pb.FrontendPutReponse, error) {
	storageArgs := StoragePut{req.Key, req.Value, 0}
	storageReply := new(string)
	funcCall := f.rpcClient.Go("Storage.StoragePut", storageArgs, &storageReply, nil)
	replyCall := <-funcCall.Done

	if replyCall.Error != nil {
		return &pb.FrontendPutReponse{
			OpId: req.OpId,
			StorageFail: true,
			Result: "",
		}, errors.New("FE to Strage fail")
	}
	// log.Print(*storageReply)

	// ret := kvslib.ResultStruct{}
	// ret.Result = storageReply
	// *reply = ret

	log.Printf("OpId: %d Put value %s to %s", req.OpId, *storageReply, req.Key)
	// log.Print(*reply.Result)


	return &pb.FrontendPutReponse{
		OpId: req.OpId,
		StorageFail: false,
		Result: *storageReply,
	}, nil
}

// !!! Old rpc function
// func (f *FrontEnd) HandlePut(args kvslib.KvslibPut, reply *kvslib.ResultStruct) error {
// 	storageArgs := StoragePut{args.Key, args.Value, 0}
// 	storageReply := new(string)
// 	funcCall := f.rpcClient.Go("Storage.StoragePut", storageArgs, &storageReply, nil)
// 	replyCall := <-funcCall.Done

// 	if replyCall.Error != nil {
// 		return errors.New("FE to Strage fail")
// 	}
// 	// log.Print(*storageReply)

// 	ret := kvslib.ResultStruct{}
// 	ret.Result = storageReply
// 	*reply = ret

// 	log.Printf("OpId: %d Put value %s to %s", args.OpId, args.Value, args.Key)
// 	// log.Print(*reply.Result)
// 	return nil
// }
