package main

import (
	"flag"
	"log"
	"strconv"
	"time"

	distkvs "example.org/cpsc416/a5"
	"example.org/cpsc416/a5/kvslib"
)

var config distkvs.ClientConfig

func main() {
	err := distkvs.ReadJSONConfig("config/client_config.json", &config)
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&config.ClientID, "id", config.ClientID, "Client ID, e.g. client1")
	flag.Parse()

	SimpleEdit()

	GetAll()

	// DelayPut()

	// MonotonicRead()

	// client := distkvs.NewClient(config, kvslib.NewKVS())
	// if err := client.Initialize(); err != nil {
	// 	log.Fatal(err)
	// }
	// // log.Printf("%v\n", client)
	// // log.Printf("PASS: client Initialize()")
	// defer client.Close()

	// if _, err := client.Get(123, "k0"); err != nil {
	// 	log.Print("errororor")
	// 	log.Println(err)
	// }

	// // log.Printf("PASS: client Get")
	// if _, err := client.Put(123, "k99", "editValue", 99); err != nil {
	// 	log.Print("errororor2")
	// 	log.Println(err)
	// }

	// // log.Printf("PASS: client Put")

	// log.Printf("Channel Result")
	// for i := 0; i < 2; i++ {
	// 	result := <-client.NotifyChannel
	// 	log.Print(result)
	// 	log.Println(*result.Result)
	// }
	// client.Close()
}

func GetAll() {
	client := distkvs.NewClient(config, kvslib.NewKVS())
	if err := client.Initialize(); err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		k := "k" + strconv.Itoa(i)
		client.Get(123, k)
	}

	for i := 0; i < 10; i++ {
		result := <-client.NotifyChannel
		log.Printf("Get Key %d => %s", i, *result.Result)
	}
	client.Close()
}

func SimpleEdit() {
	client := distkvs.NewClient(config, kvslib.NewKVS())
	if err := client.Initialize(); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		k := "k" + strconv.Itoa(i)
		client.Get(123, k)
	}

	for i := 0; i < 5; i++ {
		result := <-client.NotifyChannel
		log.Printf("Key %d => %s", i, *result.Result)
	}

	time.Sleep(2 * time.Second)

	for i := 0; i < 5; i++ {
		k := "k" + strconv.Itoa(i)
		v := strconv.Itoa(i + 10)
		client.Put(123, k, v, 0)
	}

	for i := 0; i < 5; i++ {
		result := <-client.NotifyChannel
		log.Printf("Put Key %d => %s", i, *result.Result)
	}

	time.Sleep(2 * time.Second)

	// for i := 0; i < 5; i++ {
	// 	k := "k" + strconv.Itoa(i)
	// 	client.Get(123, k)
	// }

	// for i := 5; i < 10; i++ {
	// 	result := <-client.NotifyChannel
	// 	log.Printf("Get Key %d => %s", i, *result.Result)
	// }

	client.Close()

}

// func DelayPut() {

// 	client1 := distkvs.NewClient(config, kvslib.NewKVS())
// 	if err := client1.Initialize(); err != nil {
// 		log.Fatal(err)
// 	}
// 	client2 := distkvs.NewClient(config, kvslib.NewKVS())
// 	if err := client2.Initialize(); err != nil {
// 		log.Fatal(err)
// 	}

// 	if _, err := client1.Get("client1", "k99"); err != nil {
// 		log.Println(err)
// 	}

// 	if _, err := client2.Put("client2", "k99", "editValue", 99); err != nil {
// 		log.Print("errororor2")
// 		log.Println(err)
// 	}

// 	if _, err := client1.Get("client1", "k99"); err != nil {
// 		log.Println(err)
// 	}

// 	client1.Get(123, "k99")
// 	client1.Get(123, "ddd")
// 	l := len(client1.NotifyChannel)
// 	for i := 0; i < l; i++ {
// 		result := <-client1.NotifyChannel
// 		log.Printf("C1 Get Key %d => %s", i, *result.Result)
// 	}

// 	l = len(client2.NotifyChannel)
// 	for i := 0; i < l; i++ {
// 		result := <-client2.NotifyChannel
// 		log.Printf("C2 Get Key %d => %s", i, *result.Result)
// 	}

// }

// func MonotonicRead() {
// 	client1 := distkvs.NewClient(config, kvslib.NewKVS())
// 	if err := client1.Initialize(); err != nil {
// 		log.Fatal(err)
// 	}
// 	client2 := distkvs.NewClient(config, kvslib.NewKVS())
// 	if err := client2.Initialize(); err != nil {
// 		log.Fatal(err)
// 	}

// 	client1.Get(123, "k99")
// 	go client2.Put("client2", "k99", "Update", 0)
// 	client1.Get(123, "k99")
// 	time.Sleep(8 * time.Second)

// 	fmt.Print("enddd")
// 	l := len(client1.NotifyChannel)
// 	for i := 0; i < l; i++ {
// 		result := <-client1.NotifyChannel
// 		log.Printf("C1 Get Key %d => %s", i, *result.Result)
// 	}

// 	l = len(client2.NotifyChannel)
// 	for i := 0; i < l; i++ {
// 		result := <-client2.NotifyChannel
// 		log.Printf("C2 Get Key %d => %s", i, *result.Result)
// 	}
// }
