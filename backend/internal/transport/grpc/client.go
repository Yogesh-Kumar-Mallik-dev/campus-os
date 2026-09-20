package grpc

import (
	"context"
	"fmt"
	"log"
	"time"

	campusv1 "github.com/Yogesh-Kumar-Mallik-dev/campus-os/backend/pkg/proto/campus/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// PersistenceClient encapsulates the gRPC connection to the TypeScript DB Layer.
type PersistenceClient struct {
	conn       *grpc.ClientConn
	dbService  campusv1.DatabaseServiceClient
	socketPath string
}

// BLOCK_GRPC_CLIENT_DIAL_001
// Purpose: Establishes a gRPC connection to the TypeScript persistence service via Unix Domain Socket.
// Inputs:  socketPath string
// Outputs: *PersistenceClient, error
// Errors:  ERR_SOCKET_DIAL_FAILED
func NewPersistenceClient(socketPath string) (*PersistenceClient, error) {
	target := "unix://" + socketPath
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("BLOCK_GRPC_CLIENT_DIAL_001: failed to dial DB layer on %s: %w", target, err)
	}

	return &PersistenceClient{
		conn:       conn,
		dbService:  campusv1.NewDatabaseServiceClient(conn),
		socketPath: socketPath,
	}, nil
}

// Close closes the underlying gRPC client connection.
func (c *PersistenceClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// BLOCK_TX_BEGIN_001
// Purpose: Initiates a stateful transaction session via the TS DB layer returning a TxToken.
func (c *PersistenceClient) BeginTx(ctx context.Context, reqID string, timeoutSec int32) (*campusv1.BeginTxResponse, error) {
	callCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.dbService.BeginTx(callCtx, &campusv1.BeginTxRequest{
		RequestId:      reqID,
		TimeoutSeconds: timeoutSec,
	})
	if err != nil {
		log.Printf("BLOCK_TX_BEGIN_001: BeginTx failed on socket %s: %v", c.socketPath, err)
		return nil, err
	}
	return resp, nil
}

// BLOCK_TX_COMMIT_001
// Purpose: Commits an active transaction session token.
func (c *PersistenceClient) CommitTx(ctx context.Context, txToken string) (*campusv1.CommitTxResponse, error) {
	callCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.dbService.CommitTx(callCtx, &campusv1.CommitTxRequest{
		TxToken: txToken,
	})
}

// BLOCK_TX_ROLLBACK_001
// Purpose: Aborts an active transaction session token.
func (c *PersistenceClient) RollbackTx(ctx context.Context, txToken string) (*campusv1.RollbackTxResponse, error) {
	callCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.dbService.RollbackTx(callCtx, &campusv1.RollbackTxRequest{
		TxToken: txToken,
	})
}
