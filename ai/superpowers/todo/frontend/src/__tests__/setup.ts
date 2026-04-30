import { afterAll, afterEach, beforeAll } from 'vitest'
import '@testing-library/react'
import { server } from './mocks/server'
import { resetMockTodos } from './mocks/handlers'

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }))
afterEach(() => {
  server.resetHandlers()
  resetMockTodos([])
})
afterAll(() => server.close())
