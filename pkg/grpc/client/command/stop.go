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

type StopCommand struct {
	client pb.WorkerServiceClient
}

func NewStopCommand(client pb.WorkerServiceClient) Runner {
	return &StopCommand{
		client: client,
	}
}

func (c *StopCommand) Run(args []string) error {
	if len(args) < 1 {
		return errors.New("you must pass an argument")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	newCtx := metadata.AppendToOutgoingContext(ctx, "jobid", args[0])
	defer cancel()
	command := pb.StopRequest{
		Id: args[0],
	}
	_, err := c.client.Stop(newCtx, &command, grpc.WaitForReady(true))
	if err != nil {
		return err
	}
	os.Stdout.WriteString(fmt.Sprintf("Job %v has been stopped\n", command.Id))
	return nil
}
