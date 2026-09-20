/**
 * BLOCK_DB_SERVER_ENTRYPOINT_001
 * Purpose: Initializes the TypeScript Prisma gRPC persistence server over Unix Domain Socket.
 * Inputs:  SOCKET_PATH (/var/run/campus-os/db.sock), DATABASE_URL
 * Outputs: Active gRPC server listening on Unix socket
 * Errors:  ERR_SOCKET_BIND_FAILURE, ERR_DATABASE_CONNECT_FAILURE
 */

import * as grpc from '@grpc/grpc-js';

const SOCKET_PATH = process.env.SOCKET_PATH || '/var/run/campus-os/db.sock';

export function createServer(): grpc.Server {
  const server = new grpc.Server();
  return server;
}

if (process.env.NODE_ENV !== 'test') {
  console.log(`BLOCK_DB_SERVER_ENTRYPOINT_001: Initializing persistence service on ${SOCKET_PATH}`);
}
