/**
 * BLOCK_DB_SERVER_MOCK_001
 * Purpose: In-memory mock persistence server over Unix Domain Socket for zero-dependency local development and mock testing.
 */

import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import * as crypto from 'crypto';
import * as fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const PROTO_ROOT = path.resolve(__dirname, '../../proto');

const SOCKET_PATH = process.env.SOCKET_PATH || '/tmp/campus-os-mock.sock';

export function createMockServer(): grpc.Server {
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

  // In-memory mock database store
  const mockUsers = new Map<string, any>();
  const mockClaimTokens = new Map<string, any>();
  let hasSuperAdminSeat = false;

  // Pre-seed a mock student
  const studentClaim = 'claim_genesis_test_demo';
  mockClaimTokens.set(studentClaim, {
    token: studentClaim,
    userId: 'u-yogesh-mock',
    username: 'yogesh.cse.2024.l',
    academicName: 'Yogesh',
    legalFullName: 'Yogesh Kumar Mallik',
    maskedPhone: '+91 98XXX-XX210',
    admissionType: 'LATERAL_ENTRY',
    entrySemester: 3,
    lateralSummary: 'Lateral Entry: Admitted to Semester 3. Prior polytechnic credits verified.',
    status: 'PENDING',
    expiresAt: new Date(Date.now() + 14 * 24 * 60 * 60 * 1000),
  });

  server.addService(authService, {
    BootstrapSuperAdmin: (call: any, cb: any) => {
      if (hasSuperAdminSeat) {
        return cb({ code: 6, message: 'Super Admin seat already occupied' });
      }
      hasSuperAdminSeat = true;
      const claimToken = `claim_genesis_${crypto.randomBytes(16).toString('hex')}`;
      cb(null, {
        claim_token: claimToken,
        ascii_qr: '████ [ MOCK ASCII QR ] ████',
        expires_at_unix: Math.floor(Date.now() / 1000) + 86400,
      });
    },

    ValidateClaimToken: (call: any, cb: any) => {
      const { claim_token } = call.request;
      const record = mockClaimTokens.get(claim_token);
      if (!record || record.status === 'CLAIMED') {
        return cb(null, { is_valid: false });
      }
      cb(null, {
        is_valid: true,
        user_id: record.userId,
        masked_phone_number: record.maskedPhone,
        username: record.username,
        academic_name: record.academicName,
        legal_full_name: record.legalFullName,
        admission_type: record.admissionType,
        entry_semester_number: record.entrySemester,
        lateral_entry_summary: record.lateralSummary,
        grace_period_remaining_seconds: 259200,
      });
    },

    VerifySIMAndSendOTP: (call: any, cb: any) => {
      const { device_carrier_phone } = call.request;
      const isMatch = Boolean(device_carrier_phone && device_carrier_phone.endsWith('43210'));
      if (!isMatch) {
        return cb(null, { sim_matched: false, otp_challenge_id: '', resend_available_in_seconds: 0 });
      }
      cb(null, {
        sim_matched: true,
        otp_challenge_id: 'otp_mock_challenge_123',
        resend_available_in_seconds: 45,
      });
    },

    VerifyOTPAndClaimAccount: (call: any, cb: any) => {
      const { claim_token, otp_code, new_password } = call.request;
      if (otp_code !== '123456') {
        return cb({ code: 3, message: 'Invalid OTP code' });
      }
      cb(null, {
        success: true,
        access_token: 'mock_jwt_access_token',
        refresh_token: 'mock_jwt_refresh_token',
        user_id: 'u-yogesh-mock',
        username: 'yogesh.cse.2024.l',
        email: 'yogesh.cse.2024.l@campus.edu',
        role_code: 'STUDENT',
      });
    },

    Login: (call: any, cb: any) => {
      const { identifier, password } = call.request;
      if (password === 'wrong') {
        return cb({ code: 16, message: 'Invalid credentials' });
      }
      cb(null, {
        access_token: 'mock_jwt_access_token',
        refresh_token: 'mock_jwt_refresh_token',
        expires_in_seconds: 900,
        user_id: 'u-mock-1',
        username: identifier,
        full_name: 'Mock Test User',
        role_codes: ['STUDENT'],
      });
    },

    RefreshToken: (call: any, cb: any) => {
      cb(null, {
        access_token: 'mock_new_access_token',
        refresh_token: 'mock_new_refresh_token',
        expires_in_seconds: 900,
      });
    },

    RevokeSession: (call: any, cb: any) => {
      cb(null, { success: true });
    },
  });

  return server;
}

if (process.env.NODE_ENV !== 'test') {
  if (fs.existsSync(SOCKET_PATH)) {
    fs.unlinkSync(SOCKET_PATH);
  }
  const server = createMockServer();
  server.bindAsync(`unix://${SOCKET_PATH}`, grpc.ServerCredentials.createInsecure(), (err) => {
    if (err) {
      console.error(`BLOCK_DB_SERVER_MOCK_ERR: Failed to bind to ${SOCKET_PATH}:`, err);
      process.exit(1);
    }
    console.log(`BLOCK_DB_SERVER_MOCK_001: Mock persistence server listening on unix://${SOCKET_PATH}`);
  });

  const shutdown = (signal: string) => {
    console.log(`\nBLOCK_DB_SERVER_MOCK_001: ${signal} received, gracefully shutting down mock server...`);
    server.tryShutdown((err) => {
      if (err) {
        server.forceShutdown();
      }
      try {
        if (fs.existsSync(SOCKET_PATH)) {
          fs.unlinkSync(SOCKET_PATH);
        }
      } catch (_) {}
      console.log('BLOCK_DB_SERVER_MOCK_001: Mock persistence server terminated cleanly.');
      process.exit(0);
    });
  };

  process.on('SIGINT', () => shutdown('SIGINT'));
  process.on('SIGTERM', () => shutdown('SIGTERM'));
}
