/**
 * BLOCK_DB_SERVER_MOCK_TEST_001
 * Co-located test suite for in-memory mock persistence server.
 */

import { describe, it, expect } from 'vitest';
import { createMockServer } from './server.mock.js';

describe('Mock Database Layer Server', () => {
  it('creates an initialized mock gRPC server instance', () => {
    const server = createMockServer();
    expect(server).toBeDefined();
  });
});
