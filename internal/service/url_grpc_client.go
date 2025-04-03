package service

import (
	"context"
	"time"

	pb "shortlink-gateway/proto"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

// URLGrpcClient implements the URLService interface using gRPC
type URLGrpcClient struct {
	client pb.URLServiceClient
	conn   *grpc.ClientConn
}

// NewURLGrpcClient creates a new URL service gRPC client
func NewURLGrpcClient(serverAddr string) (*URLGrpcClient, error) {
	// Set connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}
	// Create connection to gRPC service
	cc, err := grpc.NewClient(serverAddr, options...)
	if err != nil {
		return nil, err
	}

	// Connect to server
	cc.Connect()

	// Wait for connection to be ready or timeout
	state := cc.GetState()
	if state != connectivity.Ready {
		if !cc.WaitForStateChange(ctx, state) {
			return nil, ctx.Err()
		}
	}

	// Create gRPC client
	client := pb.NewURLServiceClient(cc)

	return &URLGrpcClient{
		client: client,
		conn:   cc,
	}, nil
}

// ShortenURL implements URLService.ShortenURL using gRPC
func (s *URLGrpcClient) ShortenURL(originalURL string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Call gRPC method
	resp, err := s.client.ShortenURL(ctx, &pb.ShortenURLRequest{
		OriginalUrl: originalURL,
	})
	if err != nil {
		return "", err
	}

	return resp.ShortId, nil
}

// ExpandURL implements URLService.ExpandURL using gRPC
func (s *URLGrpcClient) ExpandURL(shortID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Call gRPC method
	resp, err := s.client.ExpandURL(ctx, &pb.ExpandURLRequest{
		ShortId: shortID,
	})
	if err != nil {
		return "", err
	}

	return resp.OriginalUrl, nil
}

// Close closes the gRPC connection
func (s *URLGrpcClient) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
