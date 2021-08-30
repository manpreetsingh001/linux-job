package server

import (
	"context"
	pb "github.office.opendns.com/quadra/linux-job/pkg/grpc/proto"
	"github.office.opendns.com/quadra/linux-job/pkg/grpc/util"
	"github.office.opendns.com/quadra/linux-job/pkg/worker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)


type clientData map[string]map[string]bool

type WorkerServer struct {
	controller *worker.Controller
	clients    clientData
	access     *sync.RWMutex
	pb.UnimplementedWorkerServiceServer
}

func NewWorkerServer() *WorkerServer {
	c := worker.NewController()

	return &WorkerServer{
		access:     &sync.RWMutex{},
		controller: c,
		clients:    make(clientData),
	}
}


func (w *WorkerServer) Start(ctx context.Context, job *pb.StartRequest) (*pb.ProcessStartStatus, error) {

	c, err := w.controller.Start(job.Cmd, job.Args)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// TODO client authorization


	res := pb.ProcessStartStatus{
          Id: c.Id,
	}
	return &res, nil
}

func (w *WorkerServer) Status(ctx context.Context, req *pb.StatusRequest) (*pb.ProcessStatus, error) {
	c, err := w.controller.Status(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return util.ConvertToPbStatus(c), nil
}

func (w *WorkerServer) Stop(ctx context.Context, req *pb.StopRequest) (*pb.ProcessStatus, error) {
	c, err := w.controller.Stop(req.Id)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return util.ConvertToPbStatus(c), nil
}

func (w *WorkerServer) StreamLogs(req *pb.StreamLogRequest, stream pb.WorkerService_StreamLogsServer) error {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	bytes, err := w.controller.WatchOutput(ctx, req.Id)
	if err != nil {
		return status.Error(codes.NotFound, err.Error())
	}

	for out := range bytes {
		if err = stream.Send(util.ConvertToPbStream(out)); err != nil {
			return err
		}
	}
	return nil
}







