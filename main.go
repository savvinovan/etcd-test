package main

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Client struct {
	cli *clientv3.Client
}

func main() {
	fmt.Println("Starting the application...")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   exampleEndpoints(),
		DialTimeout: time.Minute,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close() // make sure to close the client

	// Create a new client
	c := &Client{cli: cli}

	done := make(chan struct{})

	go func() {
		for i := 0; i < 10; i++ {
			c.sendMessageToETCD(i)
		}
		close(done)
	}()
	go func() {
		c.recieveMessageFromETCD()
	}()
	<-done
	fmt.Println("Application finished.")
}

func exampleEndpoints() []string {
	return []string{"localhost:2379"}
}

func (c *Client) sendMessageToETCD(i int) {
	fmt.Println("Sending message to etcd...")

	_, err := c.cli.Put(context.TODO(), "foo", "bar"+fmt.Sprintf("%d", i))
	if err != nil {
		log.Fatal(err)
	}

}

func (c *Client) recieveMessageFromETCD() {
	fmt.Println("Receiving message from etcd...")

	wCh := c.cli.Watch(context.Background(), "foo")

	for wResp := range wCh {
		for _, ev := range wResp.Events {
			fmt.Printf("Received event: %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}
