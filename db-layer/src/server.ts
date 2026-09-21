/**
 * BLOCK_DB_SERVER_ENTRYPOINT_001
 * Purpose: Initializes the TypeScript Prisma gRPC persistence server over Unix Domain Socket.
 * Inputs:  SOCKET_PATH (/var/run/campus-os/db.sock), DATABASE_URL
 * Outputs: Active gRPC server listening on Unix socket
 * Errors:  ERR_SOCKET_BIND_FAILURE, ERR_DATABASE_CONNECT_FAILURE
 */

import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import { PrismaClient } from '@prisma/client';
import path from 'path';
import { fileURLToPath } from 'url';
import { AuthHandler } from './handles/auth.js';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const PROTO_ROOT = path.resolve(__dirname, '../../proto');

const SOCKET_PATH = process.env.SOCKET_PATH || '/var/run/campus-os/db.sock';

export function createServer(prisma?: PrismaClient): grpc.Server {
  const db = prisma || (process.env.NODE_ENV !== 'test' ? new PrismaClient() : {} as any);
  const server = new grpc.Server();

  const packageDefinition = protoLoader.loadSync(
    path.join(PROTO_ROOT, 'campus/v1/auth.proto'),
    {
      keepCase: true,
      longs: String,
      enums: String,
      defaults: true,
      oneofs: true,
      includeDirs: [PROTO_ROOT],
    }
  );

  const protoDescriptor = grpc.loadPackageDefinition(packageDefinition) as any;
  const authService = protoDescriptor.campus.v1.AuthService.service;

  const authHandler = new AuthHandler(db);

  server.addService(authService, {
    BootstrapSuperAdmin: (call: any, cb: any) => authHandler.bootstrapSuperAdmin(call, cb),
    GenerateExecutiveQR: (call: any, cb: any) => authHandler.generateExecutiveQR(call, cb),
    GenerateBulkStudentQR: (call: any, cb: any) => authHandler.generateBulkStudentQR(call, cb),
    ValidateClaimToken: (call: any, cb: any) => authHandler.validateClaimToken(call, cb),
    VerifySIMAndSendOTP: (call: any, cb: any) => authHandler.verifySIMAndSendOTP(call, cb),
    VerifyOTPAndClaimAccount: (call: any, cb: any) => authHandler.verifyOTPAndClaimAccount(call, cb),
    Login: (call: any, cb: any) => authHandler.login(call, cb),
    RefreshToken: (call: any, cb: any) => authHandler.refreshToken(call, cb),
    RevokeSession: (call: any, cb: any) => authHandler.revokeSession(call, cb),
  });

  return server;
}

if (process.env.NODE_ENV !== 'test') {
  console.log(`BLOCK_DB_SERVER_ENTRYPOINT_001: Initializing persistence service on ${SOCKET_PATH}`);
}
