package service

import (
	"context"
	"time"

	pb "github.com/hohotang/shortlink-gateway/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// URLGrpcClient implements the URLService interface using gRPC
type URLGrpcClient struct {
	client pb.URLServiceClient
	conn   *grpc.ClientConn
}

// NewURLGrpcClient creates a new URL service gRPC client
func NewURLGrpcClient(serverAddr string) (*URLGrpcClient, error) {
	// 設置連接超時
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 創建到gRPC服務的連接
	conn, err := grpc.DialContext(
		ctx,
		serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	// 創建gRPC客戶端
	client := pb.NewURLServiceClient(conn)

	return &URLGrpcClient{
		client: client,
		conn:   conn,
	}, nil
}

// ShortenURL implements URLService.ShortenURL using gRPC
func (s *URLGrpcClient) ShortenURL(originalURL string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 調用gRPC方法
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

	// 調用gRPC方法
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
