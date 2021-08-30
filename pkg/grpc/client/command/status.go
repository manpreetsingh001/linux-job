package command

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/metadata"
	"os"
	"time"

	pb "github.office.opendns.com/quadra/linux-job/pkg/grpc/proto"
	"google.golang.org/grpc"
)

type StatusCommand struct {
	client pb.WorkerServiceClient
}

func NewStatusCommand(client pb.WorkerServiceClient) Runner {
	return &StatusCommand{
		client: client,
	}
}

func (c *StatusCommand) Run(args []string) error {
	if len(args) < 1 {
		return errors.New("you must pass an argument")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	newCtx := metadata.AppendToOutgoingContext(ctx, "jobid", args[0])
	defer cancel()
	command := pb.StatusRequest{
		Id: args[0],
	}

	res, err := c.client.Status(newCtx, &command, grpc.WaitForReady(true))
	if err != nil {
		return err
	}
	os.Stdout.WriteString(fmt.Sprintf("Pid: %v Exit code: %v", res.PID, res.ExitCode))
	printStatus(res)
	return nil
}

func printStatus(s *pb.ProcessStatus) {
	var state string

	switch s.State {
	case pb.JobState_Fatal:
		state = "fatal"
	case pb.JobState_Finished:
		state = "finished"
	default:
		state = "running"
	}

	fmt.Printf("Job ID: %15s\n", s.Id)
	fmt.Printf("Job Status: %15s\n", state)
	fmt.Printf("Start time: %15s\n", time.Unix(0, s.StartTime))
}
