package command

import (
	"errors"
	"fmt"
	conf "github.office.opendns.com/quadra/linux-job/config"
	"github.office.opendns.com/quadra/linux-job/pkg/grpc/client"
	pb "github.office.opendns.com/quadra/linux-job/pkg/grpc/proto"
	"time"
)

type Runner interface {
	// Run runs a initialized runner.
	Run(args []string) error
}

func Execute(config conf.Config, args []string) error {
	if len(args) < 1 {
		return errors.New("you must pass a command")
	}
	client, err := client.NewWorkerClient(config)
	if err != nil {
		return err
	}
	cmds := map[string]Runner{
		"start":  NewStartCommand(client),
		"status":  NewStatusCommand(client),
		"stop":   NewStopCommand(client),
		"stream": NewStreamCommand(client),
	}
	cmd, ok := cmds[args[0]]
	if ok {
		return cmd.Run(args[1:])
	}
	return fmt.Errorf("unknown command: %s", cmd)
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