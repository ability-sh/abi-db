package client

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestPut(t *testing.T) {

	conn, err := grpc.Dial("127.0.0.1:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		t.Fatal(err)
		return
	}

	client := NewClient(conn)

	ctx := context.Background()

	err = client.Put(ctx, "test", []byte("OK"))

	if err != nil {
		t.Fatal(err)
		return
	}
}

func TestGet(t *testing.T) {

	conn, err := grpc.Dial("127.0.0.1:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		t.Fatal(err)
		return
	}

	client := NewClient(conn)

	ctx := context.Background()

	rs, err := client.Get(ctx, "test")

	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log(string(rs))
}

func TestMergeObject(t *testing.T) {

	conn, err := grpc.Dial("127.0.0.1:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		t.Fatal(err)
		return
	}

	client := NewClient(conn)

	ctx := context.Background()

	err = client.MergeObject(ctx, "object", map[string]interface{}{"title": "OK"})

	if err != nil {
		t.Fatal(err)
		return
	}

	err = client.MergeObject(ctx, "object", map[string]interface{}{"ctime": time.Now().Unix()})

	if err != nil {
		t.Fatal(err)
		return
	}

}

func TestGetObject(t *testing.T) {

	conn, err := grpc.Dial("127.0.0.1:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		t.Fatal(err)
		return
	}

	client := NewClient(conn)

	ctx := context.Background()

	rs, err := client.GetObject(ctx, "object")

	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log(rs)
}
