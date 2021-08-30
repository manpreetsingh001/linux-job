package command

import (
	"bufio"
	"context"
	"errors"
	"github.com/oleiade/lane"
	"google.golang.org/grpc/metadata"
	"io"
	"sync"

	pb "github.office.opendns.com/quadra/linux-job/pkg/grpc/proto"
	"google.golang.org/grpc"
)

type StreamCommand struct {
	client pb.WorkerServiceClient
}

func NewStreamCommand(client pb.WorkerServiceClient) Runner {
	return &StreamCommand{
		client: client,
	}
}

func (c *StreamCommand) Run(args []string) error {
	if len(args) < 1 {
		return errors.New("you must pass an argument")
	}
	ctx, cancel := context.WithCancel(context.Background())
	newCtx := metadata.AppendToOutgoingContext(ctx, "jobid", args[0])
	command := pb.StreamLogRequest{
		Id: args[0],
	}
	res, err := c.client.StreamLogs(newCtx, &command, grpc.WaitForReady(true))
	if err != nil {
		cancel()
		return err
	}

	ch := make(chan OutputMessage, 2)

	go readOutput(ctx, res, ch)
	return nil
}


type byteBuf struct {
	buf []byte
	err error
}

type streamReader struct {
	ch   chan byteBuf
	bufs *lane.Queue
	err  error
}

type Line struct {
	Msg     string
	Channel pb.StreamChannel
}

type OutputMessage struct {
	Err    error
	Output Line
}

func newStreamReader(ch chan byteBuf) *streamReader {
	return &streamReader{
		ch:   ch,
		bufs: lane.NewQueue(),
	}
}

func readOutput(ctx context.Context, stream pb.WorkerService_StreamLogsClient, ch chan OutputMessage) {
	defer func() {
		close(ch)
	}()

	stdOut := make(chan byteBuf)
	stdErr := make(chan byteBuf)
	stdOutRdr := newStreamReader(stdOut)
	stdErrRdr := newStreamReader(stdErr)
	errCh := make(chan error)
	wg := &sync.WaitGroup{}

	go func() {
		err := <-errCh
		stdOut <- byteBuf{err: err}
		stdErr <- byteBuf{err: err}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdOutRdr)

		for scanner.Scan() {
			line := Line{Msg: scanner.Text(), Channel: pb.StreamChannel_StdOut}
			ch <- OutputMessage{Output: line}
		}

		if err := scanner.Err(); err != nil {
			ch <- OutputMessage{Err: err}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdErrRdr)

		for scanner.Scan() {
			line := Line{Msg: scanner.Text(), Channel: pb.StreamChannel_StdErr}
			ch <- OutputMessage{Output: line}
		}

		if err := scanner.Err(); err != nil {
			ch <- OutputMessage{Err: err}
		}
	}()

line:
	for {
		select {
		case <-ctx.Done():
			errCh <- io.EOF
			break line
		default:
		}

		out, err := stream.Recv()

		if err != nil {
			errCh <- err
			break
		}

		switch out.Channel {
		case pb.StreamChannel_StdOut:
			stdOut <- byteBuf{buf: out.Output}
		default:
			stdErr <- byteBuf{buf: out.Output}
		}
	}

	wg.Wait()
}

func (s *streamReader) Read(b []byte) (int, error) {
	if s.err == io.EOF && s.bufs.Empty() {
		return 0, io.EOF
	} else if s.err != nil {
		return 0, s.err
	}

	read := 0

	for {
		select {
		case buf := <-s.ch:
			if buf.err != nil {
				s.err = buf.err
			} else {
				s.bufs.Enqueue(&(buf.buf))
			}
		default:
		}

		for s.bufs.Head() != nil && read < len(b) {
			buf := s.bufs.Head().(*[]byte)
			n := copy(b[read:], (*buf)[0:])

			if len(*buf) > len(b[read:]) {
				*buf = (*buf)[n:]
			} else {
				s.bufs.Pop()
			}

			read += n
		}

		if read > 0 || s.err != nil {
			break
		}
	}

	return read, nil

}