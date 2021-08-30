package command

import (
	"context"
	"errors"
	"fmt"
	pb "github.office.opendns.com/quadra/linux-job/pkg/grpc/proto"
	"google.golang.org/grpc"
	"os"
	"time"
)

type StartCommand struct {
	client pb.WorkerServiceClient
}

func NewStartCommand(client pb.WorkerServiceClient) Runner {
	return &StartCommand{
		client: client,
	}
}

func (c *StartCommand) Run(args []string) error {
	if len(args) < 1 {
		return errors.New("you must pass a program name")
	}
	var cargs []string
	if len(args) > 1 {
		cargs = append(cargs, args[1:]...)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	command := pb.StartRequest{
		Cmd: args[0],
		Args: cargs,
	}
	res, err := c.client.Start(ctx, &command, grpc.WaitForReady(true))
	if err != nil {
		return err
	}
	os.Stdout.WriteString(fmt.Sprintf("Job %v is started\n", res.Id))
	return nil
}

