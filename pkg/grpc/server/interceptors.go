package server

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)


func (w *WorkerServer) UnaryAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	newCtx, err := w.authorize(ctx, info.FullMethod)
	if err != nil {
		return nil, err
	}
	return handler(newCtx, req)
}


func (w *WorkerServer) StreamAuthInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	_, err := w.authorize(stream.Context(), info.FullMethod)
	if err != nil {
		return err
	}
	return handler(srv, stream)
}

func (w *WorkerServer) clientAccessToJob(client, jobId string) bool {
	w.access.RLock()
	defer w.access.RUnlock()
	jobs, ok := w.clients[client]
	if !ok {
		return false
	}
	_, ok = jobs[jobId]
	return ok
}


func (w *WorkerServer) authorize(ctx context.Context, method string) (context.Context, error) {
	// reads the peer information from the context
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("error to read peer info")
	}
	// reads user tls inforation
	tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return nil, errors.New("error to get auth info")
	}
	certs := tlsInfo.State.VerifiedChains
	if len(certs) == 0 || len(certs[0]) == 0 {
		return nil, errors.New("missing certificate chain")
	}


	cn := tlsInfo.State.VerifiedChains[0][0].Subject.CommonName

	if method == "/pb.WorkerService/Start" {
		newCtx := context.WithValue(ctx, "clientID", cn)
		return newCtx, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.NotFound, "missing metadata")
	}
	jobId, ok := md["jobid"]
	if !ok {
		return nil, status.Error(codes.NotFound, "missing jobid in metadata")
	}

	if w.clientAccessToJob(cn, jobId[0]) {
		return nil, status.Error(codes.PermissionDenied, "Not authorized to access the job")
	}
	return nil, nil
}

