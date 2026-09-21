package grpc

import (
	"context"
	"errors"
	"net"
	"path/filepath"
	"testing"

	campusv1 "github.com/Yogesh-Kumar-Mallik-dev/campus-os/backend/pkg/proto/campus/v1"
	"google.golang.org/grpc"
)

type mockDBServer struct {
	campusv1.UnimplementedDatabaseServiceServer
	beginTxFunc    func(ctx context.Context, req *campusv1.BeginTxRequest) (*campusv1.BeginTxResponse, error)
	commitTxFunc   func(ctx context.Context, req *campusv1.CommitTxRequest) (*campusv1.CommitTxResponse, error)
	rollbackTxFunc func(ctx context.Context, req *campusv1.RollbackTxRequest) (*campusv1.RollbackTxResponse, error)
}

func (m *mockDBServer) BeginTx(ctx context.Context, req *campusv1.BeginTxRequest) (*campusv1.BeginTxResponse, error) {
	if m.beginTxFunc != nil {
		return m.beginTxFunc(ctx, req)
	}
	return &campusv1.BeginTxResponse{
		TxToken:       "tx_mock_123",
		ExpiresAtUnix: 1726910000,
	}, nil
}

func (m *mockDBServer) CommitTx(ctx context.Context, req *campusv1.CommitTxRequest) (*campusv1.CommitTxResponse, error) {
	if m.commitTxFunc != nil {
		return m.commitTxFunc(ctx, req)
	}
	return &campusv1.CommitTxResponse{Success: true}, nil
}

func (m *mockDBServer) RollbackTx(ctx context.Context, req *campusv1.RollbackTxRequest) (*campusv1.RollbackTxResponse, error) {
	if m.rollbackTxFunc != nil {
		return m.rollbackTxFunc(ctx, req)
	}
	return &campusv1.RollbackTxResponse{Success: true}, nil
}

type mockAuthServer struct {
	campusv1.UnimplementedAuthServiceServer
}

func startMockGRPCServer(t *testing.T, dbSrv *mockDBServer) (string, func()) {
	t.Helper()
	sockDir := t.TempDir()
	sockPath := filepath.Join(sockDir, "test.sock")

	listener, err := net.Listen("unix", sockPath)
	if err != nil {
		t.Fatalf("failed to listen on unix socket %s: %v", sockPath, err)
	}

	server := grpc.NewServer()
	campusv1.RegisterDatabaseServiceServer(server, dbSrv)
	campusv1.RegisterAuthServiceServer(server, &mockAuthServer{})

	go func() {
		_ = server.Serve(listener)
	}()

	cleanup := func() {
		server.Stop()
		_ = listener.Close()
	}

	return sockPath, cleanup
}

func TestPersistenceClient_Lifecycle(t *testing.T) {
	mockDB := &mockDBServer{}
	sockPath, cleanup := startMockGRPCServer(t, mockDB)
	defer cleanup()

	client, err := NewPersistenceClient(sockPath)
	if err != nil {
		t.Fatalf("failed to create PersistenceClient: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Errorf("failed to close client: %v", err)
		}
	}()

	if client.AuthService() == nil {
		t.Fatal("expected AuthService() to return non-nil client")
	}

	// Test BeginTx success
	ctx := context.Background()
	beginResp, err := client.BeginTx(ctx, "req-101", 30)
	if err != nil {
		t.Fatalf("BeginTx failed: %v", err)
	}
	if beginResp.TxToken != "tx_mock_123" {
		t.Errorf("expected tx_mock_123, got %s", beginResp.TxToken)
	}

	// Test CommitTx success
	commitResp, err := client.CommitTx(ctx, beginResp.TxToken)
	if err != nil {
		t.Fatalf("CommitTx failed: %v", err)
	}
	if !commitResp.Success {
		t.Errorf("expected commit Success to be true")
	}

	// Test RollbackTx success
	rollbackResp, err := client.RollbackTx(ctx, beginResp.TxToken)
	if err != nil {
		t.Fatalf("RollbackTx failed: %v", err)
	}
	if !rollbackResp.Success {
		t.Errorf("expected rollback Success to be true")
	}
}

func TestPersistenceClient_BeginTxError(t *testing.T) {
	mockDB := &mockDBServer{
		beginTxFunc: func(ctx context.Context, req *campusv1.BeginTxRequest) (*campusv1.BeginTxResponse, error) {
			return nil, errors.New("database locked")
		},
	}
	sockPath, cleanup := startMockGRPCServer(t, mockDB)
	defer cleanup()

	client, err := NewPersistenceClient(sockPath)
	if err != nil {
		t.Fatalf("failed to create PersistenceClient: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	_, err = client.BeginTx(ctx, "req-err", 10)
	if err == nil {
		t.Fatal("expected BeginTx to fail, got nil")
	}
}
