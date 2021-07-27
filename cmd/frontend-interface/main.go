package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"net/rpc"

	"example.org/cpsc416/a5/kvslib"
	"example.org/cpsc416/a5/pb"
	"google.golang.org/grpc"
)

type FrontendInterface struct {
	rpcClient *rpc.Client
}

func main() {
	frontendAddr := flag.String("feaddr", ":8000", "address of frontend service")

	flag.Parse()

	intf := FrontendInterface{}
	err := intf.Start(*frontendAddr)
	if err != nil {
		log.Fatal(err)
	}
}

func (f *FrontendInterface) Start(frontendAddr string) error {
	// start new grpc server
	server := new(FrontendInterface)

	grpcServer := grpc.NewServer()

	pb.RegisterFrontendServer(grpcServer, server)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	// Dial to frontend rpc server
	rpcClient, err := rpc.DialHTTP("tcp", frontendAddr)
	if err != nil {
		return err 
	}
	// Set frontend rpc client
	f.rpcClient = rpcClient

	log.Println("Serving on port :50051")
	func() {
		log.Fatal(grpcServer.Serve(lis))
	}()

	return errors.New("start frontend interface failed")
}

func (f *FrontendInterface) HandleGet(ctx context.Context, req *pb.FrontendGetRequest) (*pb.FrontendGetReponse, error) {
	// Convert from grpc request to rpc request
	args := kvslib.KvslibGet{
		ClientId: req.ClientId,
		OpId: req.OpId,
		Key: req.Key,
	}
	log.Println("args", args)
	reply := new(kvslib.ResultStruct)
	// Call rpc request to Frontend and get rpc response
	funcCall := f.rpcClient.Go("FrontEnd.HandleGet", args, &reply, nil)
	log.Println(funcCall)
	replyCall := <- funcCall.Done

	if replyCall.Error != nil {
		return &pb.FrontendGetReponse{}, errors.New("key not found")
	}

	// Return nil if error
	// Convert rpc response to grpc response
	log.Println("reply", reply)
	res := &pb.FrontendGetReponse{
		OpId: req.OpId,
		StorageFail: reply.StorageFail,
		Result: *reply.Result,
	}
	log.Println("converted!", res)
	// Return grpc response to client

	return res, nil
}

func (f *FrontendInterface) HandlePut(ctx context.Context, req *pb.FrontendPutRequest) (*pb.FrontendPutReponse, error) {
	args := kvslib.KvslibPut{
		ClientId: req.ClientId,
		OpId: req.OpId,
		Key: req.Key,
		Value: req.Value,
		Delay: int(req.Delay),
	}
	reply := new(kvslib.ResultStruct)
	funcCall := f.rpcClient.Go("FrontEnd.HandlePut", args, &reply, nil)
	replyCall := <- funcCall.Done

	if replyCall.Error != nil {
		return nil, errors.New("key not found")
	}

	res := &pb.FrontendPutReponse{
		OpId: req.OpId,
		StorageFail: reply.StorageFail,
		Result: *reply.Result,
	}

	return res, nil
}