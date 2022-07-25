package srv

import (
	"context"

	"github.com/ability-sh/abi-db/pb"
	"github.com/ability-sh/abi-db/source"
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

	b, err := ss.db.Collection(task.Collection).Get(dbc, task.Key)

	if err != nil {
		if err == source.ErrNoSuchKey {
			return &pb.GetResult{Errno: 404, Errmsg: err.Error()}, nil
		}
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

	err = ss.db.Collection(task.Collection).Put(dbc, task.Key, task.Data)

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

	err = ss.db.Collection(task.Collection).MergeObject(dbc, task.Key, object)

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

	err = ss.db.Collection(task.Collection).Del(dbc, task.Key)

	if err != nil {
		return &pb.DelResult{Errno: 500, Errmsg: err.Error()}, nil
	}

	return &pb.DelResult{Errno: 200}, nil

}

func (s *server) Exec(c context.Context, task *pb.ExecTask) (*pb.ExecResult, error) {

	ctx := micro.GetContext(c)

	ss, err := GetDBService(ctx, SERVICE_ABI_DB)

	if err != nil {
		return &pb.ExecResult{Errno: 500, Errmsg: err.Error()}, nil
	}

	dbc := ss.db.NewContext()

	defer dbc.Recycle()

	data, err := ss.db.Collection(task.Collection).Exec(dbc, task.Code)

	if err != nil {
		return &pb.ExecResult{Errno: 500, Errmsg: err.Error()}, nil
	}

	return &pb.ExecResult{Errno: 200, Data: data}, nil
}

func Reg(s *grpc.Server) {
	pb.RegisterServiceServer(s, &server{})
}
