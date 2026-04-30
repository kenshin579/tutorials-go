import { defineConfig, devices } from '@playwright/test'

const PORT_FE = 5173
const PORT_BE = 8080
const BASE_URL = `http://localhost:${PORT_FE}`

export default defineConfig({
  testDir: './e2e',
  fullyParallel: false, // 단일 in-memory store를 공유하므로 직렬 실행
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: 1,
  reporter: [['list']],
  use: {
    baseURL: BASE_URL,
    trace: 'on-first-retry',
  },
  projects: [
    { name: 'chromium', use: { ...devices['Desktop Chrome'] } },
  ],
  webServer: [
    {
      command: 'cd ../backend && go run .',
      port: PORT_BE,
      reuseExistingServer: !process.env.CI,
      timeout: 30_000,
    },
    {
      command: 'npm run dev',
      port: PORT_FE,
      reuseExistingServer: !process.env.CI,
      timeout: 30_000,
    },
  ],
})
