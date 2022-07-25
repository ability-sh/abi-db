package client

import (
	"context"

	"github.com/ability-sh/abi-db/pb"
	"github.com/ability-sh/abi-lib/dynamic"
	"github.com/ability-sh/abi-lib/errors"
	"github.com/ability-sh/abi-lib/eval"
	"github.com/ability-sh/abi-lib/json"
	"google.golang.org/grpc"
)

type Collection struct {
	cli pb.ServiceClient
	key string
}

func (c *Collection) Put(cc context.Context, key string, data []byte) error {
	rs, err := c.cli.Put(cc, &pb.PutTask{Collection: c.key, Key: key, Data: data})
	if err != nil {
		return err
	}
	if rs.Errno == 200 {
		return nil
	}
	return errors.Errorf(rs.Errno, "%s", rs.Errmsg)
}

func (c *Collection) Del(cc context.Context, key string) error {
	rs, err := c.cli.Del(cc, &pb.DelTask{Collection: c.key, Key: key})
	if err != nil {
		return err
	}
	if rs.Errno == 200 {
		return nil
	}
	return errors.Errorf(rs.Errno, "%s", rs.Errmsg)
}

func (c *Collection) Get(cc context.Context, key string) ([]byte, error) {
	rs, err := c.cli.Get(cc, &pb.GetTask{Collection: c.key, Key: key})
	if err != nil {
		return nil, err
	}
	if rs.Errno == 200 {
		return rs.Data, nil
	}
	return nil, errors.Errorf(rs.Errno, "%s", rs.Errmsg)
}

func (c *Collection) GetObject(cc context.Context, key string) (interface{}, error) {
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

func (c *Collection) PutObject(cc context.Context, key string, object interface{}) error {
	b, err := json.Marshal(object)
	if err != nil {
		return err
	}
	return c.Put(cc, key, b)
}

func (c *Collection) MergeObject(cc context.Context, key string, object interface{}) error {
	b, err := json.Marshal(object)
	if err != nil {
		return err
	}
	rs, err := c.cli.Merge(cc, &pb.MergeTask{Collection: c.key, Key: key, Value: string(b)})
	if err != nil {
		return err
	}
	if rs.Errno == 200 {
		return nil
	}
	return errors.Errorf(rs.Errno, "%s", rs.Errmsg)
}

func (c *Collection) Exec(cc context.Context, code string, data interface{}) (string, error) {
	rs, err := c.cli.Exec(cc, &pb.ExecTask{Collection: c.key, Code: eval.ParseEval(code, func(key string) string {
		b, _ := json.Marshal(dynamic.Get(data, key))
		return string(b)
	})})
	if err != nil {
		return "", err
	}
	if rs.Errno == 200 {
		return rs.Data, nil
	}
	return "", errors.Errorf(rs.Errno, "%s", rs.Errmsg)
}

type Client struct {
	cli pb.ServiceClient
}

func NewClient(conn grpc.ClientConnInterface) *Client {
	return &Client{cli: pb.NewServiceClient(conn)}
}

func (c *Client) Collection(key string) *Collection {
	return &Collection{key: key, cli: c.cli}
}
