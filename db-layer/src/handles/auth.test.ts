/**
 * BLOCK_AUTH_HANDLE_TEST_001
 * Co-located test suite for AuthHandler gRPC persistence logic.
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { AuthHandler } from './auth.js';

describe('AuthHandler', () => {
  let mockPrisma: any;
  let handler: AuthHandler;

  beforeEach(() => {
    mockPrisma = {
      superAdminSeat: {
        findFirst: vi.fn(),
        create: vi.fn(),
      },
      user: {
        findFirst: vi.fn(),
        findUnique: vi.fn(),
        create: vi.fn(),
        update: vi.fn(),
      },
      role: {
        findUnique: vi.fn(),
        create: vi.fn(),
      },
      roleAssignment: {
        create: vi.fn(),
      },
      claimToken: {
        findUnique: vi.fn(),
        create: vi.fn(),
        update: vi.fn(),
        updateMany: vi.fn(),
      },
      refreshToken: {
        findUnique: vi.fn(),
        create: vi.fn(),
        update: vi.fn(),
        updateMany: vi.fn(),
      },
      studentProfile: {
        findUnique: vi.fn(),
        findMany: vi.fn(),
      },
    };

    handler = new AuthHandler(mockPrisma);
  });

  it('BLOCK_AUTH_TEST_001: bootstrapSuperAdmin prevents duplicate superadmin seats when already active', async () => {
    mockPrisma.superAdminSeat.findFirst.mockResolvedValue({
      id: 'singleton',
      activeUserId: 'u-123',
      activeUser: { isActive: true },
    });

    const call = {
      request: {
        chairperson_name: 'Dr. Sharma',
        chairperson_email: 'sharma@trust.edu',
        chairperson_phone: '+919876543210',
      },
    };

    const callback = vi.fn();
    await handler.bootstrapSuperAdmin(call, callback);

    expect(callback).toHaveBeenCalledWith({
      code: 6, // ALREADY_EXISTS
      message: expect.stringContaining('already occupied'),
    });
  });

  it('BLOCK_AUTH_TEST_001b: bootstrapSuperAdmin reissues claim docket when seat exists but is unactivated', async () => {
    mockPrisma.superAdminSeat.findFirst.mockResolvedValue({
      id: 'singleton',
      activeUserId: 'u-123',
      activeUser: { isActive: false },
    });
    mockPrisma.claimToken.updateMany.mockResolvedValue({ count: 1 });
    mockPrisma.claimToken.create.mockResolvedValue({});

    const call = {
      request: {
        chairperson_name: 'Dr. Sharma',
        chairperson_email: 'sharma@trust.edu',
        chairperson_phone: '+919876543210',
      },
    };

    const callback = vi.fn();
    await handler.bootstrapSuperAdmin(call, callback);

    expect(callback).toHaveBeenCalledWith(null, expect.objectContaining({
      claim_token: expect.stringMatching(/^claim_genesis_/),
    }));
  });

  it('BLOCK_AUTH_TEST_002: bootstrapSuperAdmin creates user, seat, and claim token when vacant', async () => {
    mockPrisma.superAdminSeat.findFirst.mockResolvedValue(null);
    mockPrisma.role.findUnique.mockResolvedValue({ id: 'r-super', code: 'SUPER_ADMIN' });
    mockPrisma.user.create.mockResolvedValue({ id: 'u-chair', username: 'chairperson.2026' });
    mockPrisma.claimToken.create.mockResolvedValue({});

    const call = {
      request: {
        chairperson_name: 'Dr. Sharma',
        chairperson_email: 'sharma@trust.edu',
        chairperson_phone: '+919876543210',
      },
    };

    const callback = vi.fn();
    await handler.bootstrapSuperAdmin(call, callback);

    expect(mockPrisma.user.create).toHaveBeenCalled();
    expect(mockPrisma.superAdminSeat.create).toHaveBeenCalledWith({
      data: { id: 'singleton', activeUserId: 'u-chair' },
    });
    expect(callback).toHaveBeenCalledWith(null, expect.objectContaining({
      claim_token: expect.stringContaining('claim_genesis_'),
      ascii_qr: expect.any(String),
    }));
  });

  it('BLOCK_AUTH_TEST_003: validateClaimToken returns student details including lateral entry', async () => {
    mockPrisma.claimToken.findUnique.mockResolvedValue({
      id: 'ct-1',
      tokenHash: 'hash123',
      userId: 'u-student',
      status: 'PENDING',
      expiresAt: new Date(Date.now() + 1000000),
      user: {
        id: 'u-student',
        username: 'yogesh.cse.2024.l',
        fullName: 'Yogesh Kumar Mallik',
        phoneNumber: '+919876543210',
      },
    });

    mockPrisma.studentProfile.findUnique.mockResolvedValue({
      userId: 'u-student',
      academicName: 'Yogesh',
      legalFullName: 'Yogesh Kumar Mallik',
      admissionType: 'LATERAL_ENTRY',
      entrySemesterNumber: 3,
    });

    const call = { request: { claim_token: 'test_token' } };
    const callback = vi.fn();

    await handler.validateClaimToken(call, callback);

    expect(callback).toHaveBeenCalledWith(null, expect.objectContaining({
      is_valid: true,
      academic_name: 'Yogesh',
      legal_full_name: 'Yogesh Kumar Mallik',
      admission_type: 'LATERAL_ENTRY',
      entry_semester_number: 3,
      masked_phone_number: expect.stringContaining('XXX-XX'),
    }));
  });

  it('BLOCK_AUTH_TEST_004: verifySIMAndSendOTP detects SIM presence match and mismatch', async () => {
    mockPrisma.claimToken.findUnique.mockResolvedValue({
      id: 'ct-1',
      status: 'PENDING',
      user: { phoneNumber: '+919876543210' },
    });

    // Mismatch test
    const mismatchCall = {
      request: {
        claim_token: 'tok_1',
        device_carrier_phone: '+919111122222',
      },
    };
    const mismatchCb = vi.fn();
    await handler.verifySIMAndSendOTP(mismatchCall, mismatchCb);
    expect(mismatchCb).toHaveBeenCalledWith(null, {
      sim_matched: false,
      otp_challenge_id: '',
      resend_available_in_seconds: 0,
    });

    // Match test
    const matchCall = {
      request: {
        claim_token: 'tok_1',
        device_carrier_phone: '+919876543210',
      },
    };
    const matchCb = vi.fn();
    await handler.verifySIMAndSendOTP(matchCall, matchCb);
    expect(matchCb).toHaveBeenCalledWith(null, expect.objectContaining({
      sim_matched: true,
      otp_challenge_id: expect.stringContaining('otp_'),
      resend_available_in_seconds: 45,
    }));
  });

  it('BLOCK_AUTH_TEST_005: refreshToken invalidates entire family upon breach reuse', async () => {
    mockPrisma.refreshToken.findUnique.mockResolvedValue({
      id: 'rt-revoked',
      tokenFamilyId: 'fam_compromised',
      isRevoked: true,
    });

    const call = { request: { refresh_token: 'stolen_token' } };
    const callback = vi.fn();

    await handler.refreshToken(call, callback);

    expect(mockPrisma.refreshToken.updateMany).toHaveBeenCalledWith({
      where: { tokenFamilyId: 'fam_compromised' },
      data: { isRevoked: true },
    });
    expect(callback).toHaveBeenCalledWith({
      code: 16,
      message: expect.stringContaining('Token breach detected'),
    });
  });

  it('BLOCK_AUTH_TEST_006: verifyOTPAndClaimAccount allows dev mock challenges without hardware SIM match', async () => {
    mockPrisma.claimToken.findUnique.mockResolvedValue({
      id: 'ct-1',
      status: 'PENDING',
      userId: 'u-chair',
      targetRole: 'SUPER_ADMIN',
      otpChallengeId: null,
      user: {
        id: 'u-chair',
        username: 'chairperson.2026',
        email: 'chairperson@campus.edu',
      },
    });
    mockPrisma.user.update.mockResolvedValue({});
    mockPrisma.claimToken.update.mockResolvedValue({});
    mockPrisma.refreshToken.create.mockResolvedValue({});

    const call = {
      request: {
        claim_token: 'claim_genesis_mock123',
        otp_challenge_id: 'mock_otp_challenge_123',
        otp_code: '123456',
        new_password: 'new_secret_password',
      },
    };
    const callback = vi.fn();

    await handler.verifyOTPAndClaimAccount(call, callback);

    expect(callback).toHaveBeenCalledWith(null, expect.objectContaining({
      success: true,
      user_id: 'u-chair',
      role_code: 'SUPER_ADMIN',
    }));
    expect(mockPrisma.user.update).toHaveBeenCalledWith(expect.objectContaining({
      where: { id: 'u-chair' },
      data: expect.objectContaining({ isActive: true }),
    }));
  });
});
