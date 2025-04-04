package service

import (
	"context"

	"github.com/hohotang/shortlink-gateway/internal/config"
	pb "github.com/hohotang/shortlink-gateway/proto"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

// URLGrpcClient implements the URLService interface using gRPC
type URLGrpcClient struct {
	client pb.URLServiceClient
	conn   *grpc.ClientConn
	cfg    *config.Config
}

// NewURLGrpcClient creates a new URL service gRPC client
func NewURLGrpcClient(serverAddr string, cfg *config.Config) (*URLGrpcClient, error) {
	// Set connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GrpcTimeout)
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
		cfg:    cfg,
	}, nil
}

// ShortenURL implements URLService.ShortenURL using gRPC
func (s *URLGrpcClient) ShortenURL(ctx context.Context, originalURL string) (string, error) {
	// Add timeout from config
	ctx, cancel := context.WithTimeout(ctx, s.cfg.GrpcTimeout)
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
func (s *URLGrpcClient) ExpandURL(ctx context.Context, shortID string) (string, error) {
	// Add timeout from config
	ctx, cancel := context.WithTimeout(ctx, s.cfg.GrpcTimeout)
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
