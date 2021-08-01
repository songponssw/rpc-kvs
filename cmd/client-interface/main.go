package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"example.org/cpsc416/a5/kvslib"
	"example.org/cpsc416/a5/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type ClientInterface struct {
	grpcClientConn *grpc.ClientConn
}

func main() {
	frontendAPIAddr := flag.String("feaddr", ":50051", "address of frontend interface server")
	clientListenAddr := flag.String("caddr", ":8080", "address of frontend interface server")

	flag.Parse()

	intf := ClientInterface{}
	err := intf.Start(*clientListenAddr, *frontendAPIAddr)

	if err != nil {
		log.Fatal(err)
	}
}


func (c *ClientInterface) Start(clientListenAddr string, frontendAPIAddr string) error {
	result := new(ClientInterface)
	rpc.Register(result)

	rpc.HandleHTTP()

	// connect to grpc 
	conn, err := grpc.Dial(frontendAPIAddr, grpc.WithInsecure(), grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Timeout: 5 * time.Second,
		PermitWithoutStream: true,
	}))
	if err != nil {
		return errors.New("connect with grpc server failed")
	}

	c.grpcClientConn = conn

	lis, err := net.Listen("tcp", clientListenAddr)
	if err != nil {
		log.Fatal("listen failed:", err)
		return err
	}

	log.Println("Serving on", clientListenAddr)
	err = http.Serve(lis, nil)
	if err != nil {
		log.Fatal("serving falied:", err)
		return err
	}

	return errors.New("client fail")
}

func (c *ClientInterface) Get(args kvslib.KvslibGet, reply *kvslib.ResultStruct) error {
	req := &pb.FrontendGetRequest{
		ClientId: "sdklfjlsdkfj",
		OpId: 44,
		Key: "sdfjhdjfhsdfAAA",
	}

	log.Println(c.grpcClientConn)

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	client := pb.NewFrontendClient(c.grpcClientConn)

	res, err := client.HandleGet(ctx, req)
	if err != nil {
		return errors.New("grpc request failed")
	}

	ret := kvslib.ResultStruct{}
	ret.Result = &res.Result
	*reply = ret

	log.Printf("response %v", res)

	log.Printf("OpId: %d Get value %s from %s", args.OpId, *reply.Result, args.Key)

	return nil
}

func (c *ClientInterface) Put(args kvslib.KvslibPut, reply *kvslib.ResultStruct) error {
	req := &pb.FrontendPutRequest{
		ClientId: args.ClientId,
		OpId: args.OpId,
		Key: args.Key,
		Value: args.Value,
		Delay: uint32(args.Delay),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	client := pb.NewFrontendClient(c.grpcClientConn)

	res, err := client.HandlePut(ctx, req)
	if err != nil {
		return errors.New("grpc request failed")
	}

	ret := kvslib.ResultStruct{}
	ret.Result = &res.Result
	*reply = ret

	log.Printf("response %v", res)

	log.Printf("OpId: %d Put value %s to %s", args.OpId, args.Value, args.Key)

	return nil
}