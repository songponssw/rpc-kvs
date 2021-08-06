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
	result := new(FrontEnd)

	// Dial to Storage service
	rpcClient, rpcErr := rpc.DialHTTP("tcp", storageAPIListenAddr)
	if rpcErr != nil {
		return errors.New("cannot established connection with Storage server")
	}

	// Assign Storage Client to struct variable
	result.rpcClient = rpcClient

	// Init gRPC server
	grpcServer := grpc.NewServer()
	// Register Frontend struct to the server
	pb.RegisterFrontendServer(grpcServer, result)

	// Listen to port
	lis, err := net.Listen("tcp", clientAPIListenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %s", clientAPIListenAddr)
	}

	// Start gRPC server (Blocking)
	log.Println("serving on port:", clientAPIListenAddr)
	log.Fatal(grpcServer.Serve(lis))

	// If code not blocking, mean Fail to serve the server
	return errors.New("frontend Fail")
}

/** 
 * TODO: Handle Client Timeout => Cancel the action 
 * TODO: Handle Retry Request on Failed
 */
// Function define by Protobuf file
func (f *FrontEnd) HandleGet(
	ctx context.Context, 
	req *pb.FrontendGetRequest,
) (*pb.FrontendGetReponse, error) {

	// convert to storage args
	storageArgs := StorageGet{req.Key}	// Storage Request Struct
	storageReply := new(StorageGetResult)	// Storage Response Struct
	funcCall := f.rpcClient.Go("Storage.StorageGet", storageArgs, &storageReply, nil) // Call Storage Func
	replyCall := <-funcCall.Done // Wait for response from Storage
	if replyCall.Error != nil {
		return nil, errors.New("FE to Storage fail")
	}

	// DEBUG: Print Value of Key gotten from Storage 
	log.Printf("OpId: %d Get value %s from %s", req.OpId, *storageReply.Value, req.Key)

	// Return Result to Protobuf
	return &pb.FrontendGetReponse{
		OpId: req.OpId,
		StorageFail: false,
		Result: *storageReply.Value,
	}, nil
}

/** 
 * TODO: Handle Client Timeout => Cancel the action 
 * TODO: Handle Retry Request on Failed
 */
// Function define by Protobuf file
func (f *FrontEnd) HandlePut(
	ctx context.Context, 
	req *pb.FrontendPutRequest, 
	) (*pb.FrontendPutReponse, error) {
	// Same as Get
	storageArgs := StoragePut{req.Key, req.Value, 0}
	storageReply := new(string)
	funcCall := f.rpcClient.Go("Storage.StoragePut", storageArgs, &storageReply, nil)
	replyCall := <-funcCall.Done
	if replyCall.Error != nil {
		// Return Empty Struct w/ Error
		return &pb.FrontendPutReponse{
			OpId: req.OpId,
			StorageFail: true,
			Result: "",
		}, errors.New("FE to Storage fail")
	}

	// DEBUG: Print Value
	log.Printf("OpId: %d Put value %s to %s", req.OpId, *storageReply, req.Key)

	// Return Result to Protobuf
	return &pb.FrontendPutReponse{
		OpId: req.OpId,
		StorageFail: false,
		Result: *storageReply,
	}, nil
}