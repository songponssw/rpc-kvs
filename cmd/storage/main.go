package main

import (
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	distkvs "example.org/cpsc416/a5"
)

func main() {
	var config distkvs.StorageConfig
	err := distkvs.ReadJSONConfig("config/storage_config.json", &config)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(config)

	storage := Storage{}

	err = storage.Start(config.FrontEndAddr, string(config.StorageAdd), config.DiskPath)
	if err != nil {
		log.Fatal(err)
	}
}

type StorageConfig struct {
	StorageID        string
	StorageAdd       distkvs.StorageAddr
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
	delay int
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

func (*Storage) Start(frontEndAddr string, storageAddr string, diskPath string) error {
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

	if _, err := os.Stat("mem"); os.IsNotExist(err) {
		_, err := os.Create("mem")
		if err != nil {
			panic(err)
		}
	}

	data := make([]byte, 10000000)
	file, err := os.Open("mem")
	if err != nil {
		log.Println(err)
	}

	s := "key not found"

	count, err := file.Read(data)
	if err != io.EOF {
		log.Fatal(err)
	} else if err == nil {
		k := ""
		v := ""
		i := 0
		state := 0

		for i < count {
			if state == 0 {
				if data[i] == ';' {
					if k == args.Key {
						state = 1
						v = ""
					} else {
						state = 2
					}
				} else {
					k = k + string(data[i])
				}
			} else if state == 1 {
				if data[i] == '\n' {
					state = 0
					i++
					s = v
				} else {
					v = v + string(data[i])
				}
			} else if state == 2 {
				if data[i] == '\n' {
					state = 0
					i++
				}
			}
			i++
		}
	}

	ret.Value = &s
	*reply = ret
	log.Printf("Get value %s from %s", *reply.Value, args.Key)

	return nil
}

func (*Storage) StoragePut(args StoragePut, reply *string) error {
	log.Printf("Put value %s to %s", args.Value, args.Key)
	s := "success"

	if args.Key == "k99" {
		log.Print("delay for 5 second")
		time.Sleep(5 * time.Second)
	} else {
		if _, err := os.Stat("mem"); os.IsNotExist(err) {
			_, err := os.Create("mem")
			if err != nil {
				panic(err)
			}
		}
		file, err := os.OpenFile("mem", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
			s = "cannot open file"
		}
		defer file.Close()
		if _, err := file.WriteString(args.Key + ";" + args.Value); err != nil {
			log.Fatal(err)
			s = "cannot write"
		}
	}

	*reply = s

	return nil
}
