/**
 * Manual QA Test Configuration
 * 
 * All credentials and environment variables for testing.
 * Update these values for different environments.
 */
export const TEST_CONFIG = {
  /** Admin panel URL */
  adminUrl: "http://100.67.206.65:5173",

  /** Server API URL */
  apiUrl: "http://localhost:8080",

  /** Test user credentials */
  credentials: {
    username: "sma",
    password: "sma",
  },

  /** WebSocket URL template */
  wsUrl: "ws://100.67.206.65:8080/ws",

  /** Server log path */
  serverLog: "/tmp/herbst-web.log",

  /** Default world context used by admin panel */
  defaultWorld: "herbst-mud",

  /** Timeout for API calls (ms) */
  apiTimeout: 10000,

  /** Timeout for page loads (ms) */
  pageTimeout: 15000,
} as const;

export type TestConfig = typeof TEST_CONFIG;
