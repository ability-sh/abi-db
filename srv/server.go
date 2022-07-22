package srv

import (
	"context"

	"github.com/ability-sh/abi-db/pb"
	"github.com/ability-sh/abi-lib/json"
	"github.com/ability-sh/abi-micro/micro"
	"google.golang.org/grpc"
)

type server struct {
}

func (s *server) Get(c context.Context, task *pb.GetTask) (*pb.GetResult, error) {

	ctx := micro.GetContext(c)

	ss, err := GetDBService(ctx, SERVICE_ABI_DB)

	if err != nil {
		return &pb.GetResult{Errno: 500, Errmsg: err.Error()}, nil
	}

	dbc := ss.db.NewContext()

	defer dbc.Recycle()

	b, err := ss.db.Get(dbc, task.Key)

	if err != nil {
		return &pb.GetResult{Errno: 500, Errmsg: err.Error()}, nil
	}

	return &pb.GetResult{Errno: 200, Data: b}, nil
}

func (s *server) Put(c context.Context, task *pb.PutTask) (*pb.PutResult, error) {

	ctx := micro.GetContext(c)

	ss, err := GetDBService(ctx, SERVICE_ABI_DB)

	if err != nil {
		return &pb.PutResult{Errno: 500, Errmsg: err.Error()}, nil
	}

	dbc := ss.db.NewContext()

	defer dbc.Recycle()

	err = ss.db.Put(dbc, task.Key, task.Data)

	if err != nil {
		return &pb.PutResult{Errno: 500, Errmsg: err.Error()}, nil
	}

	return &pb.PutResult{Errno: 200}, nil

}

func (s *server) Merge(c context.Context, task *pb.MergeTask) (*pb.MergeResult, error) {

	ctx := micro.GetContext(c)

	ss, err := GetDBService(ctx, SERVICE_ABI_DB)

	if err != nil {
		return &pb.MergeResult{Errno: 500, Errmsg: err.Error()}, nil
	}

	dbc := ss.db.NewContext()

	defer dbc.Recycle()

	var object interface{} = nil

	err = json.Unmarshal([]byte(task.Value), &object)

	if err != nil {
		return &pb.MergeResult{Errno: 500, Errmsg: err.Error()}, nil
	}

	err = ss.db.MergeObject(dbc, task.Key, object)

	if err != nil {
		return &pb.MergeResult{Errno: 500, Errmsg: err.Error()}, nil
	}

	return &pb.MergeResult{Errno: 200}, nil
}

func (s *server) Del(c context.Context, task *pb.DelTask) (*pb.DelResult, error) {

	ctx := micro.GetContext(c)

	ss, err := GetDBService(ctx, SERVICE_ABI_DB)

	if err != nil {
		return &pb.DelResult{Errno: 500, Errmsg: err.Error()}, nil
	}

	dbc := ss.db.NewContext()

	defer dbc.Recycle()

	err = ss.db.Del(dbc, task.Key)

	if err != nil {
		return &pb.DelResult{Errno: 500, Errmsg: err.Error()}, nil
	}

	return &pb.DelResult{Errno: 200}, nil

}

func Reg(s *grpc.Server) {
	pb.RegisterServiceServer(s, &server{})
}
