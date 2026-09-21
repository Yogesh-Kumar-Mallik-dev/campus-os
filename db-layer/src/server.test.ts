import { describe, it, expect } from 'vitest';
import { createServer } from './server.js';

/**
 * BLOCK_DB_SERVER_TEST_001
 * Purpose: Verifies gRPC persistence server creation.
 */
describe('Database Layer Server', () => {
  it('creates an initialized gRPC server instance', () => {
    const mockPrisma = {} as any;
    const server = createServer(mockPrisma);
    expect(server).toBeDefined();
  });
});
