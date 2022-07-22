package client

import (
	"context"

	"github.com/ability-sh/abi-db/pb"
	"github.com/ability-sh/abi-lib/json"
	"google.golang.org/grpc"
)

type Error struct {
	Errno  int32  `json:"errno"`
	Errmsg string `json:"errmsg"`
}

func (E *Error) Error() string {
	return E.Errmsg
}

type Client struct {
	cli pb.ServiceClient
}

func NewClient(conn grpc.ClientConnInterface) *Client {
	return &Client{cli: pb.NewServiceClient(conn)}
}

func (c *Client) Put(cc context.Context, key string, data []byte) error {
	rs, err := c.cli.Put(cc, &pb.PutTask{Key: key, Data: data})
	if err != nil {
		return err
	}
	if rs.Errno == 200 {
		return nil
	}
	return &Error{Errno: rs.Errno, Errmsg: rs.Errmsg}
}

func (c *Client) Del(cc context.Context, key string) error {
	rs, err := c.cli.Del(cc, &pb.DelTask{Key: key})
	if err != nil {
		return err
	}
	if rs.Errno == 200 {
		return nil
	}
	return &Error{Errno: rs.Errno, Errmsg: rs.Errmsg}
}

func (c *Client) Get(cc context.Context, key string) ([]byte, error) {
	rs, err := c.cli.Get(cc, &pb.GetTask{Key: key})
	if err != nil {
		return nil, err
	}
	if rs.Errno == 200 {
		return rs.Data, nil
	}
	return nil, &Error{Errno: rs.Errno, Errmsg: rs.Errmsg}
}

func (c *Client) GetObject(cc context.Context, key string) (interface{}, error) {
	b, err := c.Get(cc, key)
	if err != nil {
		return nil, err
	}

	var object interface{} = nil

	err = json.Unmarshal(b, &object)

	if err != nil {
		return nil, err
	}

	return object, nil
}

func (c *Client) PutObject(cc context.Context, key string, object interface{}) error {
	b, err := json.Marshal(object)
	if err != nil {
		return err
	}
	return c.Put(cc, key, b)
}

func (c *Client) MergeObject(cc context.Context, key string, object interface{}) error {
	b, err := json.Marshal(object)
	if err != nil {
		return err
	}
	rs, err := c.cli.Merge(cc, &pb.MergeTask{Key: key, Value: string(b)})
	if err != nil {
		return err
	}
	if rs.Errno == 200 {
		return nil
	}
	return &Error{Errno: rs.Errno, Errmsg: rs.Errmsg}
}
