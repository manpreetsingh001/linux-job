package command

import (
"context"
"errors"
"fmt"
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
	defer cancel()
	command := pb.StatusRequest{
		Id: args[0],
	}
	res, err := c.client.Status(ctx, &command, grpc.WaitForReady(true))
	if err != nil {
		return err
	}
	os.Stdout.WriteString(fmt.Sprintf("Pid: %v Exit code: %v", res.PID, res.ExitCode))
	return nil
}

