package util

import (
	pb "github.office.opendns.com/quadra/linux-job/pkg/grpc/proto"
	"github.office.opendns.com/quadra/linux-job/pkg/worker"
	"github.office.opendns.com/quadra/linux-job/pkg/worker/exec"
	"github.office.opendns.com/quadra/linux-job/pkg/worker/log"
)

func ConvertToPbStatus(p *worker.ProcessStatus) *pb.ProcessStatus {
	var s pb.JobState

	switch p.Status.State {
	case exec.Finished:
		s = pb.JobState_Finished
	case exec.Fatal:
		s = pb.JobState_Fatal
	default:
		s = pb.JobState_Running
	}

	err := ""
	if p.Status.Error != nil {
		err = p.Status.Error.Error()
	}

	return &pb.ProcessStatus{
		Id:        p.Id,
		PID:       int64(p.Status.PID),
		Cmd:       p.Status.Cmd,
		State:     s,
		ExitCode:  int64(p.Status.ExitCode),
		Error:     err,
		StartTime: p.Status.StartTime.UnixNano(),
		StopTime:  p.Status.StopTime.UnixNano(),
	}
}

func ConvertToPbStream(l log.LogOutput) *pb.ProcessStream {
	var c pb.StreamChannel
	switch l.Type {
	case worker.StdOut:
		c = pb.StreamChannel_StdOut
	default:
		c = pb.StreamChannel_StdErr
	}
	return &pb.ProcessStream{Output: l.Bytes, Channel: c}
}
